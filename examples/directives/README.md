# Environment File Directives

This directory contains examples demonstrating the use of directives in `.env` files. Directives are special instructions that control how environment variables are processed.

## Available Directives

### `#remove` Directive

Removes specified environment variables from the existing set before merging.

**Syntax:** `#remove KEY1 KEY2 KEY3...`

**Example:**
```env
#remove OLD_KEY DEPRECATED_KEY
NEW_KEY=new_value
```

### `#require` Directive

Fails processing if specified environment variables are not present in the final merged output.

**Syntax:** `#require KEY1 KEY2 KEY3...`

**Example:**
```env
#require DATABASE_URL API_KEY
SOME_KEY=some_value
```

### `#filter` Directive

Removes environment variables based on key names or wildcard patterns from the final merged output.

**Syntax:** `#filter PATTERN1 PATTERN2 PATTERN3...`

**Pattern Types:**
- **Exact match:** `KEY_NAME` - removes variables with exact key names (case-insensitive)
- **Wildcard patterns:** `*` matches any sequence of characters
  - `TEST_*` - removes all keys starting with "TEST_"
  - `*_PROD` - removes all keys ending with "_PROD"
  - `API_*_KEY` - removes keys matching the pattern "API_" + anything + "_KEY"

**Examples:**
```env
# Remove specific keys
#filter DEBUG_KEY LOG_LEVEL

# Remove all test-related keys
#filter TEST_* *_TEST

# Remove all production keys
#filter *_PROD PROD_*

# Remove API keys with specific pattern
#filter API_*_KEY *_API_*
```

## Directive Processing Order

Directives are processed in the following order:

1. **`#remove`** - Applied to existing variables before merging
2. **Variable merging** - File variables override existing ones
3. **`#filter`** - Applied to the merged result
4. **`#require`** - Applied to the final result (fails if required variables are missing)

## Combining Directives

You can use multiple directives in the same file:

```env
# Remove old variables
#remove OLD_KEY DEPRECATED_KEY

# Filter out test and debug variables
#filter TEST_* DEBUG_* *_DEV

# Require essential variables
#require DATABASE_URL API_KEY SECRET_TOKEN

# Add new variables
NEW_KEY=new_value
FINAL_KEY=final_value
```

## Case Insensitivity

All directive names are case-insensitive:
- `#remove` === `#REMOVE` === `#Remove`
- `#require` === `#REQUIRE` === `#Require`
- `#filter` === `#FILTER` === `#Filter`

## Examples

See the example files in this directory:
- `base.env` - Basic environment variables
- `production.env` - Production configuration with directives
