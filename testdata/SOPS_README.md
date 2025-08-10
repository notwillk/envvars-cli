# SOPS Test Files

This directory contains test files for SOPS (Secrets OPerationS) encryption testing.

## Test Files

### 1. `secrets.yaml`
A simple configuration file with database and API secrets.

### 2. `production.yaml`
A production configuration file with nested service configurations and security settings.

## Creating Encrypted Files

To create encrypted versions of these files for testing, you'll need to install SOPS and use age encryption.

### Install SOPS
```bash
# On macOS
brew install sops

# On Linux
wget -O sops https://github.com/getsops/sops/releases/download/v3.10.2/sops-v3.10.2.linux
chmod +x sops
sudo mv sops /usr/local/bin/
```

### Install age
```bash
# On macOS
brew install age

# On Linux
sudo apt-get install age
```

### Generate age key pair
```bash
age-keygen -o test-key.txt
```

This will create:
- `test-key.txt` - Your private key (keep this secret)
- A public key that looks like: `age1ql3z7hjy54pw3hyww5ay3fgkd...`

### Encrypt the files
```bash
# Encrypt secrets.yaml
sops --encrypt --age age1ql3z7hjy54pw3hyww5ay3fgkd... secrets.yaml > secrets.enc.yaml

# Encrypt production.yaml
sops --encrypt --age age1ql3z7hjy54pw3hyww5ay3fgkd... production.yaml > production.enc.yaml
```

## Test Decryption Keys

For testing purposes, use these decryption keys:

### For `secrets.enc.yaml`:
```
age1ql3z7hjy54pw3hyww5ay3fgkd...
```

### For `production.enc.yaml`:
```
age1ql3z7hjy54pw3hyww5ay3fgkd...
```

## Testing with the CLI

Once you have encrypted files, test them with:

```bash
# Test SOPS decryption
go run main.go --sops "age1ql3z7hjy54pw3hyww5ay3fgkd...@testdata/secrets.enc.yaml" --format json

# Test with multiple sources
go run main.go \
  --env testdata/basic.env \
  --sops "age1ql3z7hjy54pw3hyww5ay3fgkd...@testdata/production.enc.yaml" \
  --format yaml

# Test verbose output
go run main.go \
  --sops "age1ql3z7hjy54pw3hyww5ay3fgkd...@testdata/secrets.enc.yaml" \
  --verbose
```

## Expected Output

The SOPS processor will:
1. Decrypt the encrypted YAML files
2. Flatten nested structures into environment variable format
3. Convert keys to uppercase with underscores
4. Handle arrays by joining with commas
5. Convert booleans and numbers to strings

For example, `services.web.secret_key` becomes `SERVICES_WEB_SECRET_KEY`.

## Security Note

These are test files with dummy secrets. Never use real secrets in test files or commit encrypted files with real keys to version control.
