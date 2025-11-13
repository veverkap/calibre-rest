#!/usr/bin/env python3
"""
extract_calibredb_commands.py
Parses the Sphinx-generated calibredb manual HTML and extracts commands
and options into structured JSON.
"""
import os
from bs4 import BeautifulSoup
import json
import re
import requests


def parse_option(option_element):
    """Parse a single option (dl.std.option) element."""
    option_data = {
        "names": [],
        "description": ""
    }
    
    # Get option names from the dt element
    dt = option_element.find("dt", class_="sig")
    if dt:
        for sig_name in dt.find_all("span", class_="sig-name"):
            name = sig_name.get_text(strip=True)
            if name:
                option_data["names"].append(name)
    
    # Get description from dd element
    dd = option_element.find("dd")
    if dd:
        # Get all text, preserving some structure
        description = dd.get_text(separator=" ", strip=True)
        option_data["description"] = description
    
    return option_data


def parse_command(command_section):
    """Parse a command section and extract all relevant information."""
    command_id = command_section.get("id", "")

    # Skip non-command sections
    if command_id in ["calibredb", "global-options", "adding-from-folders", "epub-options"]:
        return None
    
    command_data = {
        "name": command_id.replace("-", "_"),
        "id": command_id,
        "description": "",
        "usage": "",
        "options": []
    }
    
    # Get command title/heading
    h2 = command_section.find("h2")
    if h2:
        title = h2.get_text(strip=True)
        # Remove the backref symbol
        title = re.sub(r'¶$', '', title)
        command_data["title"] = title
    
    # Get usage/syntax (typically in a code block)
    code_block = command_section.find("div", class_="highlight-none")
    if code_block:
        pre = code_block.find("pre")
        if pre:
            command_data["usage"] = pre.get_text(strip=True)
    
    # Get description (first paragraph after usage)
    for p in command_section.find_all("p", recursive=False):
        desc = p.get_text(separator=" ", strip=True)
        if desc and not desc.startswith("Whenever you pass arguments"):
            command_data["description"] = desc
            break
    
    # Parse all options for this command
    for dl_option in command_section.find_all("dl", class_="std option"):
        option = parse_option(dl_option)
        if option["names"]:
            command_data["options"].append(option)
    
    # Also check subsections for options (like "Adding From Folders", "Epub Options")
    for subsection in command_section.find_all("section", recursive=False):
        subsection_name = subsection.get("id", "")
        for dl_option in subsection.find_all("dl", class_="std option"):
            option = parse_option(dl_option)
            if option["names"]:
                # Add subsection context to the option
                option["subsection"] = subsection_name
                command_data["options"].append(option)
    
    # deduplicate options
    seen_options = set()
    unique_options = []
    for option in command_data["options"]:
        option_key = tuple(option["names"])
        if option_key not in seen_options:
            seen_options.add(option_key)
            unique_options.append(option)
    command_data["options"] = unique_options
    
    return command_data


def parse_global_options(soup):
    """Parse the global options section."""
    global_section = soup.find("section", id="global-options")
    if not global_section:
        return []
    
    global_options = []
    for dl_option in global_section.find_all("dl", class_="std option"):
        option = parse_option(dl_option)
        if option["names"]:
            global_options.append(option)
    
    return global_options


def main():
    """Main function to parse HTML and output JSON."""
    # Read HTML from file or stdin
    URL = "https://manual.calibre-ebook.com/generated/en/calibredb.html"
    print(f"Fetching {URL} ...")
    r = requests.get(URL, timeout=30)
    r.raise_for_status()    
    html_content = r.text
    soup = BeautifulSoup(html_content, 'html.parser')
    
    # Find the main calibredb section
    main_section = soup.find("section", id="calibredb")
    if not main_section:
        print("Error: Could not find main calibredb section")
        return
    
    # Parse global options
    print("Parsing global options...")
    global_options = parse_global_options(main_section)
    
    # Parse all command sections
    print("Parsing commands...")
    commands = []
    for section in main_section.find_all("section", recursive=False):
        section_id = section.get("id", "")
        if section_id == "global-options":
            continue
        
        command_data = parse_command(section)

        if command_data:
            commands.append(command_data)

    # sort commands by name
    commands.sort(key=lambda x: x["name"])
    # Build final output
    output = {
        "global_options": global_options,
        "commands": commands
    }
    
    # Write to JSON file
    output_file = "scraped.json"
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(output, f, indent=2, ensure_ascii=False)
    
    print(f"\n✓ Successfully extracted {len(commands)} commands")
    print(f"✓ Output written to: {output_file}")


if __name__ == "__main__":
    main()