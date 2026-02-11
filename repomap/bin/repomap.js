#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

// detect platform and arch
const platform = os.platform();
const arch = os.arch();

let binaryName = '';

if (platform === 'linux') {
    if (arch === 'x64') binaryName = 'repomap-linux-amd64';
    else if (arch === 'arm64') binaryName = 'repomap-linux-arm64';
} else if (platform === 'darwin') {
    if (arch === 'x64') binaryName = 'repomap-darwin-amd64';
    else if (arch === 'arm64') binaryName = 'repomap-darwin-arm64';
} else if (platform === 'win32') {
    if (arch === 'x64') binaryName = 'repomap-windows-amd64.exe';
}

if (!binaryName) {
    console.error(`Unsupported platform: ${platform}-${arch}`);
    process.exit(1);
}

// Path to the compiled binary in build/ directory
const binPath = path.join(__dirname, '..', 'build', binaryName);

// args to pass through
const args = process.argv.slice(2);

const child = spawn(binPath, args, {
    stdio: 'inherit'
});

child.on('exit', (code) => {
    process.exit(code);
});

child.on('error', (err) => {
    console.error(`Failed to start repomap binary: ${err.message}`);
    process.exit(1);
});
