#!/usr/bin/env bash
set -euo pipefail

FOLDER_URL="https://api.github.com/repos/kovidgoyal/calibre/contents/src/calibre/db/cli"
OUTPUT_JSON="$(pwd)/calibredb_cli_options.json"

echo "==> Fetching Python files from $FOLDER_URL ..."
mkdir -p calibre_cli
cd calibre_cli

# Download all files from the GitHub API folder listing
files=$(curl -s "$FOLDER_URL" | jq -r '.[] | select(.type=="file") | .download_url')

if [ -z "$files" ]; then
  echo "❌ No files found. Check network or jq installation."
  exit 1
fi

echo "==> Downloading $(echo "$files" | wc -l) Python files..."
for url in $files; do
    fname=$(basename "$url")
    echo "  - $fname"
    curl -s -O "$url"
done
echo "✓ Download complete."


# --- Python AST extractor ---
cat << 'PYCODE' > extract_options.py
import ast, json
from pathlib import Path
from re import A
import sys


def eval_ast(node):
    """Evaluate string-like AST nodes, supporting _('...') and .format()."""
    if isinstance(node, ast.Constant):
        return node.value
    if isinstance(node, ast.Str):  # Py<3.8
        return node.s

    # Handle _('...') wrapper
    if isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == "_":
        if node.args:
            return eval_ast(node.args[0])

    # Handle _("...").format(...)
    if isinstance(node, ast.Call) and isinstance(node.func, ast.Attribute) and node.func.attr == "format":
        base = eval_ast(node.func.value)
        args = [eval_ast(a) for a in node.args]
        try:
            return base.format(*args)
        except Exception:
            return f"{base} (format args: {args})"

    try:
        return ast.unparse(node)
    except Exception:
        return None


def parse_cli_file(py_file: Path):
    """Parse a calibre CLI command file for parser options and groups."""
    tree = ast.parse(py_file.read_text(encoding="utf-8"))
    command = {"file": py_file.name, "usage_text": None, "options": [], "option_groups": []}

    # Detect get_parser() usage/help text
    for node in ast.walk(tree):
        if isinstance(node, ast.Call) and isinstance(node.func, ast.Name) and node.func.id == "get_parser":
            if node.args:
                text = eval_ast(node.args[0])
                # check if text includes 'join('
                if "join(" in text:
                    print(f"⚠️  Warning: usage text in {py_file.name} includes 'join(', may be complex.")
                command["usage_text"] = text
            break

    # Keep track of group variables
    groups = {}

    for node in ast.walk(tree):
        # Detect parser.add_option_group()
        if (
            isinstance(node, ast.Call)
            and hasattr(node.func, "attr")
            and node.func.attr == "add_option_group"
        ):
            group_name = None
            if node.args:
                group_name = eval_ast(node.args[0])
            target = None
            # find assignment to variable: group = parser.add_option_group(...)
            parent = getattr(node, "parent", None)
            if isinstance(parent, ast.Assign) and len(parent.targets) == 1:
                target = parent.targets[0]
                if isinstance(target, ast.Name):
                    groups[target.id] = {"name": group_name, "options": []}
            continue

        # Detect add_option() calls
        if (
            isinstance(node, ast.Call)
            and hasattr(node.func, "attr")
            and node.func.attr == "add_option"
        ):
            opt = {"_file": py_file.name, "_lineno": node.lineno, "_node": ast.unparse(node)}
            for kw in node.keywords:
                key = kw.arg
                opt[key] = eval_ast(kw.value)
            flags = [eval_ast(a) for a in node.args if isinstance(a, (ast.Str, ast.Constant))]
            if flags:
                opt["flags"] = flags

            # Determine if it belongs to a group
            if isinstance(node.func.value, ast.Name) and node.func.value.id in groups:
                groups[node.func.value.id]["options"].append(opt)
            else:
                command["options"].append(opt)

    # Merge groups
    command["option_groups"] = list(groups.values())
    return command


def attach_parents(tree):
    """Annotate each AST node with a .parent attribute."""
    for node in ast.walk(tree):
        for child in ast.iter_child_nodes(node):
            child.parent = node

CLI_DIR = Path(".")
all_commands = {}

for py_file in sorted(CLI_DIR.glob("cmd*.py")):
    print(f"  • Parsing {py_file.name}")
    try:
        parsed = parse_cli_file(py_file)
        all_commands[py_file.stem] = parsed
        print(
            f"    → usage: {bool(parsed['usage_text'])}, "
            f"options: {len(parsed['options'])}, "
            f"groups: {len(parsed['option_groups'])}"
        )
    except Exception as e:
        print(f"    ⚠️  Parse error in {py_file.name}: {e}")

print(f"==> Parsed {len(all_commands)} CLI command files.")

output_path = Path("calibredb_cli_options.json")
with open(output_path, "w", encoding="utf-8") as f:
    json.dump(all_commands, f, indent=2, ensure_ascii=False)
print(f"✓ Wrote {output_path.resolve()}")

# --- End of Python AST extractor ---
PYCODE

echo "==> Running Python extractor ..."
python3 extract_options.py

cp calibredb_cli_options.json "$OUTPUT_JSON"
echo "✓ Output written to: $OUTPUT_JSON"

echo "==> Cleaning up downloaded Python files ..."
cd ..
rm -rf calibre_cli
echo "✓ Done."