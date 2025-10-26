#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}📦 JotterXpress RPM Package Builder${NC}"
echo ""

# Check if running as root (for some operations)
if [ "$EUID" -eq 0 ]; then 
    echo -e "${YELLOW}⚠ Warning: Running as root. Some operations may require sudo.${NC}"
fi

# Variables
APP_NAME="jotterxpress"
VERSION="1.0.0"
TARBALL="${APP_NAME}-${VERSION}.tar.gz"
SPEC_FILE="${APP_NAME}.spec"
BUILD_DIR=$(pwd)
RPM_BUILD_DIR="${HOME}/rpmbuild"

echo "🔨 Building application..."
go build -o jotterxpress cmd/jotterxpress/main.go

if [ ! -f "jotterxpress" ]; then
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Build successful!${NC}"
echo ""

# Create source tarball
echo "📦 Creating source tarball..."
tar --exclude='.git' \
    --exclude='rpmbuild' \
    --exclude='*.tar.gz' \
    --exclude='jotterxpress.spec' \
    --exclude='package-rpm.sh' \
    --exclude='bin' \
    --exclude='.vscode' \
    --exclude='*.md' \
    -czf "${TARBALL}" .

if [ ! -f "$TARBALL" ]; then
    echo -e "${RED}❌ Failed to create tarball!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Tarball created: ${TARBALL}${NC}"
echo ""

# Setup RPM build directory structure
echo "📁 Setting up RPM build directories..."
mkdir -p ${RPM_BUILD_DIR}/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}

# Copy files to RPM build directories
cp "${TARBALL}" ${RPM_BUILD_DIR}/SOURCES/
cp "${SPEC_FILE}" ${RPM_BUILD_DIR}/SPECS/

# Build RPM
echo "🔨 Building RPM package..."
rpmbuild -ba ${RPM_BUILD_DIR}/SPECS/${SPEC_FILE}

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ RPM package built successfully!${NC}"
    echo ""
    echo "📦 Package location:"
    echo "   Binary RPM: ${RPM_BUILD_DIR}/RPMS/*/jotterxpress-${VERSION}-*.rpm"
    echo "   Source RPM: ${RPM_BUILD_DIR}/SRPMS/jotterxpress-${VERSION}-*.src.rpm"
    echo ""
    echo "🚀 To install the package, run:"
    echo "   sudo rpm -ivh ${RPM_BUILD_DIR}/RPMS/*/jotterxpress-${VERSION}-*.rpm"
    echo ""
    echo "🔍 To verify installation:"
    echo "   rpm -qlp ${RPM_BUILD_DIR}/RPMS/*/jotterxpress-${VERSION}-*.rpm"
else
    echo -e "${RED}❌ RPM build failed!${NC}"
    exit 1
fi
