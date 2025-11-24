# PDF report dependency

The PDF/email features use `wkhtmltopdf`. When the API runs outside Docker you must install it locally so the binary is on `PATH` (or point `WKHTMLTOPDF_PATH` to its directory).

- Debian/Ubuntu: `sudo apt-get update && sudo apt-get install -y wkhtmltopdf`
- macOS (Homebrew): `brew install wkhtmltopdf`
- Windows (Chocolatey): `choco install wkhtmltopdf` then add `C:\Program Files\wkhtmltopdf\bin` to `PATH`.

Check with `wkhtmltopdf --version` and restart the API. The Docker image used by `docker-compose` already installs this dependency.
