# Legacy-Activator 🚀

A Go-based cross-platform application designed to activate older iOS devices by leveraging various exploits and a custom backend.

## Badges 🛡️

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Description 📝

Legacy-Activator is a tool developed in Go that aims to activate older iOS devices. It utilizes exploits tailored for specific iOS versions and relies on a backend server (written in PHP) to serve necessary configuration files (`.plist`). The application provides a graphical user interface (GUI) built with Fyne for ease of use.

It supports activation for a range of iOS versions, utilizing different exploitation methods (e.g., `itunesstored` for newer versions, SFTP for older versions) based on the detected device and its iOS firmware.

## Table of Contents 📚

- [Project Title & Badges](#legacy-activator-🚀)
- [Description](#description-📝)
- [Table of Contents](#table-of-contents-📚)
- [Features](#features-⭐)
- [Tech Stack](#tech-stack-💻)
- [Installation](#installation-⬇️)
- [Usage](#usage-💡)
- [Project Structure](#project-structure-📁)
- [API Reference](#api-reference-🌐)
- [Contributing](#contributing-🤝)
- [License](#license-⚖️)
- [Important Links](#important-links-🔗)
- [Footer](#footer-✨)

## Features ⭐

- **Cross-Platform GUI Application:** Built with Fyne for a native desktop experience.
- **iOS Version Support:** Supports activation for a wide range of iOS versions (e.g., iOS 4.0 through 10.3.4).
- **Multiple Exploitation Methods:** Implements different strategies based on iOS version:
  - `itunesstored` exploit for iOS 8.0 - 10.3.4.
  - SFTP patch for iOS 4.0 - 8.4.1.
- **Custom Backend Integration:** Communicates with a PHP-based backend server to fetch device-specific activation data (`.plist` files).
- **Device Detection:** Automatically detects connected iOS devices via `usbmuxd`.
- **Configurable Settings:** Allows users to configure backend URL and retry attempts.
- **Status Updates:** Provides real-time status updates during the activation process.

## Tech Stack 💻

- **Languages:** Go, PHP, Markdown
- **Frameworks/Libraries:**
  - **Go:** Fyne (GUI), `go-resty/resty` (HTTP client), `ncruces/go-sqlite3` (SQLite), `go-pkg/ios` (iOS communication), `golang.org/x/crypto/ssh` (SSH client).
  - **PHP:** (Backend server logic)
- **Other:** `usbmuxd` (for device communication)

## Installation ⬇️

**Prerequisites:**

* **Go Toolchain:** Ensure you have Go installed (version 1.26.5 or compatible).
* **Libimobiledevice:** Ensure `libimobiledevice` is installed and functional on your system for `usbmuxd` communication.
* **PHP Backend:** You need to set up the PHP backend server provided in the `backend/` directory.

**Steps:**

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/Chrissie21/Legacy-Activator.git
   cd Legacy-Activator
   ```

2. **Set up the PHP Backend:**
   - Navigate to the `backend/` directory.
   - Ensure your web server (e.g., Apache, Nginx) is configured to serve the `server.php` file.
   - Place the `patched.plist` files in the appropriate `backend/plists/[model]/[build]/` subdirectories.
   - The default backend URL in the Go application is `http://localhost/backend/server.php`. Adjust this if your backend is hosted elsewhere.

3. **Build the Go Application:**
   ```bash
   go build -o legacy-activator ./cmd/api/
   ```

4. **Run the Application:**
   ```bash
   ./legacy-activator
   ```

## Usage 💡

Legacy-Activator is designed to activate older iOS devices, particularly those with older iOS versions that may no longer be actively supported or have limited activation options.

**Use Cases:**

* **Activating older iPhones/iPads:** Use this tool to activate devices running iOS versions between 4.0 and 10.3.4.
* **Bypassing activation lock (if applicable):** This tool might be used in scenarios where a device needs activation to be functional.

**How to Use:**

1. **Connect your iOS device** to your computer.
2. **Ensure the device is recognized** by your system (e.g., through `libimobiledevice` tools or by appearing in Finder/iTunes).
3. **Launch the `legacy-activator` application** (after building it).
4. The application will attempt to detect your connected device.
5. If a compatible device is found and the backend is reachable, the **"Activate Device"** button will become enabled.
6. Click **"Activate Device"** to initiate the activation process.
7. Follow any on-screen prompts. The application will display status updates.
8. **Settings:** Click the gear icon (⚙) to configure the **Backend URL** and the number of **Attempts** for the activation process.

**Important:**
* Ensure your device has a stable internet connection (Wi-Fi) during the activation process, especially for newer supported iOS versions.
* For older iOS versions, the process might involve rebooting the device and patching system files.

## Project Structure 📁

```
Legacy-Activator/
├── backend/
│   ├── plists/
│   │   └── [device_model]/
│   │       └── [build_version]/
│   │           └── patched.plist
│   └── server.php
├── cmd/api/
│   └── main.go
├── pkg/
│   ├── exploits/
│   │   ├── common/
│   │   │   └── common.go
│   │   ├── itunesstored/
│   │   │   ├── itunesstored.go
│   │   │   └── payload.sql
│   │   └── sftp/
│   │       ├── patch.go
│   │       └── sftp.go
├── .DS_Store
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## API Reference 🌐

The `backend/server.php` script exposes two main functionalities:

1.  **Support Query Endpoint:**
    *   **URL:** `[backend_url]?support_query&model=[device_model]&build=[build_version]`
    *   **Method:** `GET`
    *   **Description:** Checks if the server supports a given device model and build version.
    *   **Response:** JSON object with a `supported` boolean field (e.g., `{"supported": true}`).

2.  **File Download Endpoint:**
    *   **URL:** Serviced implicitly when a device connects and its `User-Agent` string is parsed.
    *   **Method:** `GET`
    *   **Description:** Serves the appropriate `patched.plist` file based on the device's User-Agent string (containing model and build information).
    *   **Response:** A `patched.plist` file or a `403 Forbidden` error.

## Contributing 🤝

Contributions are welcome! Please feel free to:

- Fork the repository.
- Create a new branch (`git checkout -b feature/your-feature-name`).
- Make your changes.
- Commit your changes (`git commit -m 'Add some feature'`).
- Push to the branch (`git push origin feature/your-feature-name`).
- Open a Pull Request.

Please ensure your code follows the existing style and includes appropriate tests if applicable.

## License ⚖️

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Important Links 🔗

- **Repository:** [https://github.com/Chrissie21/Legacy-Activator](https://github.com/Chrissie21/Legacy-Activator)

## Footer ✨

This README was generated based on the analysis of the **Legacy-Activator** repository.

- **Repository:** [Legacy-Activator](https://github.com/Chrissie21/Legacy-Activator)
- **Author:** Chrissie21
- **Contact:** (Not specified in repository)

If you find this project useful, consider starring ⭐, forking 🍴, and opening issues 🐛 for any bugs or feature requests!


---
**<p align="center">Generated by [ReadmeCodeGen](https://www.readmecodegen.com/)</p>**