# Legacy-Activator

A Go-based cross-platform application designed to activate older iOS devices by leveraging various exploits and a custom backend.

## Description

Legacy-Activator is a tool developed in Go that aims to activate older iOS devices. It utilizes exploits tailored for specific iOS versions and relies on a backend server (written in PHP) to serve necessary configuration files (`.plist`). The application provides a graphical user interface (GUI) built with **Fyne** for ease of use.

It supports activation for a range of iOS versions, utilizing different exploitation methods (e.g., `itunesstored` for newer versions, SFTP for older versions) based on the detected device and its iOS firmware.

---

## Technical Architecture & Exploit Deep-Dive 

Legacy-Activator isolates its core bypassing mechanics into two distinct system pipelines. Rather than breaking cryptographic chain-of-trust over the network, it works by falsifying local system states directly on the hardware.

### 1. The `itunesstored` Database Hijack (iOS 8.0 - 10.3.4)
For later legacy versions, the application targets the `itunesstored` daemon (iTunes Store Daemon), which processes storefront synchronization, hardware pairing tokens, and cloud account mappings.

* **The Vulnerability:** `itunesstored` maintains persistent local caches of activation tickets, unique device certificates, and authorization signatures in internal SQLite containers. In these iOS iterations, the daemon does not enforce strict cryptographic integrity checks or process isolation over its underlying databases when handling system restores or localized configuration syncing via properties.
* **Low-Level Execution Steps (`pkg/exploits/itunesstored/`):**
  1. The application pairs with the target device utilizing native USB multiplexing (`usbmuxd`) via `go-pkg/ios`, establishing a session with the device's native restore or Apple File Conduit (AFC) services.
  2. The application stages `pkg/exploits/itunesstored/payload.sql`.
  3. This script injects records containing valid wildcard parameters into tables corresponding to device activations and account parameters.
  4. By directly altering variables inside the cached sqlite layout, `itunesstored` interprets the modified environment state as fully authenticated. It reports a status of `ActivationState = 'Activated'` to the iOS lock manager, bypassing the need for a remote handshake with Apple's ALBERT activation infrastructure.

### 2. The SFTP Lockdown Filesystem Patch (iOS 4.0 - 8.4.1)
For earlier iOS architectures, activation enforcement relies on explicit property list states managed directly inside the device framework structure.

* **The Vulnerability:** The security model allows the core device handshake daemon (`lockdownd`) to read unencrypted configurations directly from the root directories if root write privileges are active.
* **Low-Level Execution Steps (`pkg/exploits/sftp/`):**
  1. This pipeline expects the device to be in a jailbroken state with an active SSH daemon running (accessible over a USB tunnel via `usbmuxd` mapped to local ports such as `2222` or `22`).
  2. Using `golang.org/x/crypto/ssh`, the Go binary opens an SFTP filesystem stream using the default operating credentials (`root:alpine`).
  3. The exploit targets the device configuration pathways, typically located at:
     ```text
     /System/Library/Lockdown/Library/ActivationRecords/
     # Or systemic caches depending on specific iOS architecture variants
     ```
  4. The tool backs up the factory activation plist and pushes your specialized `patched.plist` generated via the PHP backend. This custom plist declares explicit parameters asserting that the setup wizard constraints are completed, effectively allowing the device to drop straight to the home screen after a daemon reload or a quick system reboot (`killall lockdownd`).

---

## Features

* **Cross-Platform GUI Application:** Built with Fyne for a native desktop experience.
* **iOS Version Support:** Supports activation for a wide range of iOS versions (iOS 4.0 through 10.3.4).
* **Multiple Exploitation Methods:** Implements different strategies based on iOS version:
  * `itunesstored` exploit for iOS 8.0 - 10.3.4.
  * SFTP patch for iOS 4.0 - 8.4.1.
* **Custom Backend Integration:** Communicates with a PHP-based backend server to fetch device-specific activation data (`.plist` files).
* **Device Detection:** Automatically detects connected iOS devices via `usbmuxd`.
* **Configurable Settings:** Allows users to configure backend URL and retry attempts.
* **Status Updates:** Provides real-time status updates during the activation process.

---

## Credits & References  

* **[TA of the Checkm8 Exploit](https://www.elabforensics.com/wp-content/uploads/sites/4/2021/01/Technical-analysis-of-the-checkm8-exploit.pdf)** – In-depth documentation covering low-level iOS hardware exploitation mechanisms. 
* **[King (pgarba/King)](https://github.com/pgarba/King)** – A clean C-port implementation reference for device-level bootrom interactions.
* **[bl_sbx (hanakim3945/bl_sbx)](https://github.com/hanakim3945/bl_sbx)** – A proof-of-concept implementation for an iOS `itunesstored` and `bookassetd` sandbox escape vulnerability.


---

## ⚠️ Disclaimer & Educational Use Only

**This project is for educational, research, and proof-of-concept purposes only.** 

By downloading, cloning, or using Legacy-Activator, you agree to the following terms:

* **No Liability:** The author(s) of this project assume absolute zero liability for any data loss, bricked devices, hardware damage, or system instability resulting from the use or misuse of this tool, its backend, or its scripts.
* **Intended Use:** This tool is designed strictly as an academic exploration of older iOS security architectures, legacy daemons (`itunesstored`, `lockdownd`), and low-level system pipelines. It should only be used on legacy hardware that you own legally.
* **No Commercial Use:** This software is source-available and explicitly prohibited from being used for commercial purposes, paid repair services, or re-packaged as a paid product. 
* **Compliance:** You are solely responsible for ensuring that your use of this software complies with all applicable local, national, and international laws, regulations, and terms of service. 
