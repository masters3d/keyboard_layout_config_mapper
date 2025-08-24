#!/bin/bash

# Keyboard Configuration Download Script
# Downloads the latest configurations for all four keyboard types

set -e  # Exit on any error

echo "🚀 Starting keyboard configuration download..."

# Function to download file
download_file() {
    local dir=$1
    local file=$2
    local url=$3
    
    echo "  📥 Downloading $file to $dir/"
    
    # Create directory if it doesn't exist
    mkdir -p "$dir"
    
    # Download file
    if curl -fsSL "$url/$file" -o "$dir/$(basename $file)"; then
        echo "    ✅ Successfully downloaded $(basename $file)"
    else
        echo "    ❌ Failed to download $file"
        return 1
    fi
}

# Download Kinesis Advantage 2 configuration (reference file only)
echo ""
echo "📁 Processing configs/kinesis2..."
KINESIS_BASE="https://raw.githubusercontent.com/masters3d/supportfiles/master"
download_file "configs/kinesis2" "1_qwerty.txt" "$KINESIS_BASE"

# Download QMK ErgoDox configuration
echo ""
echo "📁 Processing configs/qmk_ergodox..."
QMK_BASE="https://raw.githubusercontent.com/masters3d/qmk_firmware/masters3d"
download_file "configs/qmk_ergodox" "keyboards/ergodox_ez/keymaps/masters3d/keymap.c" "$QMK_BASE"

# Download Glove80 configuration
echo ""
echo "📁 Processing configs/zmk_glove80..."
GLOVE80_BASE="https://raw.githubusercontent.com/masters3d/glove80-zmk-config/cheyo"
download_file "configs/zmk_glove80" "config/glove80.keymap" "$GLOVE80_BASE"

# Download Advantage360 configuration
echo ""
echo "📁 Processing configs/zmk_adv360..."
ADV360_BASE="https://raw.githubusercontent.com/masters3d/Adv360-Pro-ZMK/cheyo"
download_file "configs/zmk_adv360" "config/adv360.keymap" "$ADV360_BASE"

echo ""
echo "🎉 All configurations downloaded successfully!"
echo ""
echo "📋 Summary:"
echo "  - configs/kinesis2/: Kinesis Advantage 2 reference configuration (1_qwerty.txt)"
echo "  - configs/qmk_ergodox/: QMK ErgoDox keymap (keymap.c)"
echo "  - configs/zmk_glove80/: Glove80 ZMK keymap (glove80.keymap)"
echo "  - configs/zmk_adv360/: Advantage360 ZMK keymap (adv360.keymap)"
echo ""
echo "🔗 Repository mapping:"
echo "  - Kinesis Advantage 2: masters3d/supportfiles/master/1_qwerty.txt"
echo "  - QMK ErgoDox: masters3d/qmk_firmware/masters3d/keyboards/ergodox_ez/keymaps/masters3d/keymap.c"
echo "  - Glove80: masters3d/glove80-zmk-config/cheyo/config/glove80.keymap"
echo "  - Advantage360: masters3d/Adv360-Pro-ZMK/cheyo/config/adv360.keymap"