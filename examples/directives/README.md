# Environment File Directives

This example demonstrates the new directive functionality in the ENV processor.

## What are Directives?

Directives are special instructions that can be embedded in environment files to control how variables are processed. They use the format `#[directive-name]` and are case-insensitive.

## Available Directives

### #remove

The `#remove` directive removes specified environment variables from the existing set before merging.

**Syntax:**
```
#remove KEY1 KEY2 KEY3
```

**Example:**
```bash
# Create a base environment file
echo "BASE_KEY=base_value" > base.env
echo "REMOVE_KEY=remove_value" >> base.env
echo "KEEP_KEY=keep_value" >> base.env

# Create a file with remove directive
cat > directives.env << EOF
#remove REMOVE_KEY
#remove BASE_KEY
NEW_KEY=new_value
EOF

# Process with directives
./bin/envvars-cli --env base.env --env directives.env
```

**Output:**
```
KEEP_KEY=keep_value
NEW_KEY=new_value
```

## How It Works

1. **Directives are processed first**: Before merging variables, all directives in a file are applied to the existing variable set
2. **Case-insensitive matching**: `#remove key` will remove `KEY`, `key`, `Key`, etc.
3. **Multiple arguments**: You can specify multiple keys to remove in a single directive
4. **Redefinition works**: If you remove a key and then redefine it in the same file, the new value will be used

## Use Cases

- **Environment-specific overrides**: Remove production keys when loading development configs
- **Security**: Remove sensitive variables before merging with public configs
- **Configuration management**: Selectively include/exclude variables based on context

## Notes

- Directives only affect variables that exist in previous sources
- Variables defined in the same file as a directive will still be processed normally
- Regular comments (starting with `# `) are ignored and don't interfere with directives
- Directives are processed in the order they appear in the file
