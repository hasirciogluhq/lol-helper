#!/bin/bash

echo "🎮 LoL Helper Build Script"
echo "=========================="
echo ""

# Renk kodları
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Bağımlılıkları kontrol et
echo -e "${BLUE}📦 Bağımlılıklar kontrol ediliyor...${NC}"
go mod download
go mod tidy

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Bağımlılıklar yüklenemedi!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Bağımlılıklar hazır${NC}"
echo ""

# Build dizini oluştur
mkdir -p build

# Platform seçimi
echo "Hangi platform için build yapmak istiyorsunuz?"
echo "1) macOS (Apple Silicon - M1/M2/M3)"
echo "2) macOS (Intel)"
echo "3) Windows"
echo "4) Linux"
echo "5) Hepsi"
read -p "Seçiminiz (1-5): " choice

build_macos_arm64() {
    echo -e "${BLUE}🍎 macOS (ARM64) için build ediliyor...${NC}"
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/lol-helper-macos-arm64 main.go
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ macOS ARM64 build başarılı${NC}"
    else
        echo -e "${RED}❌ macOS ARM64 build başarısız${NC}"
    fi
}

build_macos_amd64() {
    echo -e "${BLUE}🍎 macOS (Intel) için build ediliyor...${NC}"
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/lol-helper-macos-amd64 main.go
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ macOS Intel build başarılı${NC}"
    else
        echo -e "${RED}❌ macOS Intel build başarısız${NC}"
    fi
}

build_windows() {
    echo -e "${BLUE}🪟 Windows için build ediliyor...${NC}"
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o build/lol-helper-windows.exe main.go
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Windows build başarılı${NC}"
    else
        echo -e "${RED}❌ Windows build başarısız${NC}"
    fi
}

build_linux() {
    echo -e "${BLUE}🐧 Linux için build ediliyor...${NC}"
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/lol-helper-linux main.go
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Linux build başarılı${NC}"
    else
        echo -e "${RED}❌ Linux build başarısız${NC}"
    fi
}

case $choice in
    1)
        build_macos_arm64
        ;;
    2)
        build_macos_amd64
        ;;
    3)
        build_windows
        ;;
    4)
        build_linux
        ;;
    5)
        build_macos_arm64
        echo ""
        build_macos_amd64
        echo ""
        build_windows
        echo ""
        build_linux
        ;;
    *)
        echo -e "${RED}❌ Geçersiz seçim!${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}🎉 Build işlemi tamamlandı!${NC}"
echo -e "${BLUE}📁 Build dosyaları: ./build/ klasöründe${NC}"
ls -lh build/
