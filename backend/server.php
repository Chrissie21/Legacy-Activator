<?php
// Define the directory where the device-specific plist files are stored.
$baseDir = __DIR__ . '/plists';

/**
 * Safely resolves the path to a patched.plist file based on the device model and build.
 * Includes security checks to prevent path traversal vulnerabilities.
 */
function resolve_plist($baseDir, $model, $build) {
    // Ensure both model and build are provided and do not contain '..' (path traversal attempt)
    if (!$model || !$build || strpos($model, '..') !== false || strpos($build, '..') !== false) {
        return null;
    }
    
    // Construct the expected file path in the format: plists/[model]/[build]/patched.plist
    $filePath = $baseDir . '/' . $model . '/' . $build . '/patched.plist';
    
    // Return the path only if the file actually exists on disk
    return file_exists($filePath) ? $filePath : null;
}

// --- API Endpoint: Support Query ---
// If the request URL contains '?support_query', act as a JSON API 
// to check if a specific model and build are supported by the server.
if (isset($_GET['support_query'])) {
    $filePath = resolve_plist($baseDir, $_GET['model'] ?? '', $_GET['build'] ?? '');
    header('Content-Type: application/json');
    // Output a boolean true/false depending on if the plist exists
    echo json_encode(['supported' => $filePath !== null]);
    exit();
}

// --- File Download Endpoint ---
// Retrieve the HTTP User-Agent string sent by the client making the request.
$userAgent = $_SERVER['HTTP_USER_AGENT'] ?? '';

$model = null;
$build = null;

// Extract the device model from the User-Agent string using a regular expression.
// Example match: model/iPhone10,1
if (preg_match('/model\/([a-zA-Z0-9,]+)/', $userAgent, $mMatches)) {
    $model = $mMatches[1];
}

// Extract the OS build number from the User-Agent string.
// Example match: build/18D52
if (preg_match('/build\/([a-zA-Z0-9]+)/', $userAgent, $bMatches)) {
    $build = $bMatches[1];
}

// Attempt to find the correct plist file for the extracted model and build.
$filePath = resolve_plist($baseDir, $model, $build);

// If a valid file is found, serve it to the client as a file download.
if ($filePath !== null) {
    // Set HTTP headers to trigger a file download on the client side
    header('Content-Description: File Transfer');
    header('Content-Encoding: identity');
    header('Content-Type: application/xml'); // plists are typically formatted as XML
    header('Content-Disposition: attachment; filename="patched.plist"');
    header('Content-Length: ' . filesize($filePath));
    header('Cache-Control: must-revalidate');
    header('Pragma: public');
    
    // Stream the file contents directly to the HTTP response
    readfile($filePath);
    exit();
}

// If no file was found or the User-Agent didn't contain valid data, reject the request.
http_response_code(403);
echo "Forbidden";
exit();
?>