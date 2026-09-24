## Legacy-Activator 
A Go-based cross-platform application designed to activate older iOS devices by leveraging various exploits and a custom backend.

## Description 
Legacy-Activator is a tool developed in Go that aims to activate older iOS devices. It utilizes exploits tailored for specific iOS versions and relies on a backend server (written in PHP) to serve necessary configuration files (.plist). The application provides a graphical user interface (GUI) built with Fyne for ease of use.
It supports activation for a range of iOS versions, utilizing different exploitation methods (e.g., itunesstored for newer versions, SFTP for older versions) based on the detected device and its iOS firmware.

## Technical Architecture & Exploit Deep-Dive 🔬
Legacy-Activator isolates its core bypassing mechanics into two distinct system pipelines. Rather than breaking cryptographic chain-of-trust over the network, it works by falsifying local system states directly on the hardware.

## 1. The itunesstored Database Hijack (iOS 8.0 - 10.3.4)
For later legacy versions, the application targets the itunesstored daemon (iTunes Store Daemon), which processes storefront synchronization, hardware pairing tokens, and cloud account mappings.

* 
* The Vulnerability: itunesstored maintains persistent local caches of activation tickets, unique device certificates, and authorization signatures in internal SQLite containers. In these iOS iterations, the daemon does not enforce strict cryptographic integrity checks or process isolation over its underlying databases when handling system restores or localized configuration syncing via properties.
* Low-Level Execution Steps (pkg/exploits/itunesstored/):
1. The application pairs with the target device utilizing native USB multiplexing (usbmuxd) via go-pkg/ios, establishing a session with the device's native restore or Apple File Conduit (AFC) services.
   2. The application stages pkg/exploits/itunesstored/payload.sql.
   3. This script injects records containing valid wildcard parameters into tables corresponding to device activations and account parameters.
   4. By directly altering variables inside the cached sqlite layout, itunesstored interprets the modified environment state as fully authenticated. It reports a status of ActivationState = 'Activated' to the iOS lock manager, bypassing the need for a remote handshake with Apple's ALBERT activation infrastructure.
* 

## 2. The SFTP Lockdown Filesystem Patch (iOS 4.0 - 8.4.1)
For earlier iOS architectures, activation enforcement relies on explicit property list states managed directly inside the device framework structure.

* 
* The Vulnerability: The security model allows the core device handshake daemon (lockdownd) to read unencrypted configurations directly from the root directories if root write privileges are active.
* Low-Level Execution Steps (pkg/exploits/sftp/):
1. This pipeline expects the device to be in a jailbroken state with an active SSH daemon running (accessible over a USB tunnel via usbmuxd mapped to local ports such as 2222 or 22).
   2. Using golang.org/x/crypto/ssh, the Go binary opens an SFTP filesystem stream using the default operating credentials (root:alpine).
   3. The exploit targets the device configuration pathways, typically located at:
   
   /System/Library/Lockdown/Library/ActivationRecords/# Or systemic caches depending on specific iOS architecture variants
   
   4. The tool backs up the factory activation plist and pushes your specialized patched.plist generated via the PHP backend. This custom plist declares explicit parameters asserting that the setup wizard constraints are completed, effectively allowing the device to drop straight to the home screen after a daemon reload or a quick system reboot (killall lockdownd).
* 

------------------------------
## Features

* 
* Cross-Platform GUI Application: Built with Fyne for a native desktop experience.
* iOS Version Support: Supports activation for a wide range of iOS versions (e.g., iOS 4.0 through 10.3.4).
* Multiple Exploitation Methods: Implements different strategies based on iOS version:
* itunesstored exploit for iOS 8.0 - 10.3.4.
   * SFTP patch for iOS 4.0 - 8.4.1.
* Custom Backend Integration: Communicates with a PHP-based backend server to fetch device-specific activation data (.plist files).
* Device Detection: Automatically detects connected iOS devices via usbmuxd.
* Configurable Settings: Allows users to configure backend URL and retry attempts.
* Status Updates: Provides real-time status updates during the activation process.
* 

## Tech Stack

* 
* Languages: Go, PHP, SQL
* Frameworks/Libraries:
* Go: Fyne (GUI), go-resty/resty (HTTP client), ncruces/go-sqlite3 (SQLite), go-pkg/ios (iOS communication), golang.org/x/crypto/ssh (SSH client).
   * PHP: (Backend server logic)
* Other: usbmuxd (for device communication)
* 

## Installation
Prerequisites:

* 
* Go Toolchain: Ensure you have Go installed (version 1.26.5 or compatible).
* Libimobiledevice: Ensure libimobiledevice is installed and functional on your system for usbmuxd communication.
* PHP Backend: You need to set up the PHP backend server provided in the backend/ directory.
* 

Steps:

   1. Clone the Repository:
   
   git clone https://github.com/Chrissie21/Legacy-Activator.git
   cd Legacy-Activator
   
   2. Set up the PHP Backend:
   * Navigate to the backend/ directory.
      * Ensure your web server (e.g., Apache, Nginx) is configured to serve the server.php file.
      * Place the patched.plist files in the appropriate backend/plists/[model]/[build]/ subdirectories.
      * The default backend URL in the Go application is http://localhost/backend/server.php. Adjust this if your backend is hosted elsewhere.
   3. Build the Go Application:
   
   go build -o legacy-activator ./cmd/api/
   
   4. Run the Application:
   
   ./legacy-activator
   
   
## Usage & Runtime Flow
Legacy-Activator streamlines complex command-line exploitation mechanics into a single click by analyzing the downstream device information automatically.
## Automated Sequence of Operations:

   1. Connect your device over USB and launch the legacy-activator runtime.
   2. The Go execution wrapper invokes usbmuxd protocols to discover the connected device signature, mapping out its internal variables: Device Model (e.g., iPhone3,1), Build Version, and unique UDID.
   3. The UI queries the PHP backend API to verify that a payload pairing exists for that unique model configuration profile.
   4. Triggering the Activation:
   * SFTP Flow: If an older architecture is identified, the client initializes connection pipes directly across your USB terminal port, prompts or completes filesystem authentication, injects your custom configuration file over the active root path, and changes file execution permissions (chmod 755).
      * Database Injection Flow: If a newer iOS version is processed, the framework targets local services via USB backup protocols, loads your specialized payload.sql instructions, edits the properties within the app container ecosystem, and safely breaks connections.
   5. A system prompt will instruct you to reboot your device to refresh running state variables.

------------------------------

## API Reference
The backend/server.php script exposes two main functionalities:

   1. Support Query Endpoint:
   * URL: [backend_url]?support_query&model=[device_model]&build=[build_version]
      * Method: GET
      * Description: Checks if the server supports a given device model and build version.
      * Response: JSON object with a supported boolean field (e.g., {"supported": true}).
   2. File Download Endpoint:
   * URL: Serviced implicitly when a device connects and its User-Agent string is parsed.
      * Method: GET
      * Description: Serves the appropriate patched.plist file based on the device's User-Agent string (containing model and build information).
      * Response: A patched.plist file or a 403 Forbidden error.
   
## Contributing
Contributions are welcome! Please feel free to:

* 
* Fork the repository.
* Create a new branch (git checkout -b feature/your-feature-name).
* Make your changes.
* Commit your changes (git commit -m 'Add some feature').
* Push to the branch (git push origin feature/your-feature-name).
* Open a Pull Request.
* 

Please ensure your code follows the existing style and includes appropriate tests if applicable.
