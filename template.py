import sys
import os
import textwrap

TEMPLATE = """
package commands

import (

	"example.com/m/v2/auth"
)

type {command_name}Command struct{{}}

func ({first_letter}c {command_name}Command) Name() string {{ return "{command_name}" }}

func ({first_letter}c {command_name}Command) Description() string {{ return "{command_description}" }}

func ({first_letter}c {command_name}Command) Exec(token *auth.FreshToken, args []string) error {{
    return nil
}}

func init() {{
	registerCommand({command_name}Command{{}})
}}
"""

def create_command_file(command_name, description = ""):
    file_name = f"{command_name}.go"
    file_path = os.path.join("commands", file_name)

    if os.path.exists(file_path):
        print(f"File '{file_path}' already exists.")
        return

    content = TEMPLATE.format(command_name= command_name.capitalize(), first_letter = command_name[0], command_description = description)

    with open(file_path, "w") as f:
        f.write(textwrap.dedent(content))

    print(f"Created command at '{file_path}'")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python template.py <command_name> <description> (description is optional)")
        sys.exit(1)
    description = ""
    description = sys.argv[2] if len(sys.argv) > 2 else ""
    command_name = sys.argv[1]
    create_command_file(command_name, description)
    
    