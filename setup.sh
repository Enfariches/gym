# Подготовка
echo "Обновление пакетов системы..."
sudo apt update && sudo apt upgrade -y

# Установка базовых утилит
echo "Установка базовых инструментов (включая make)..."
sudo apt install -y ca-certificates curl gnupg make

# Установка Docker
echo "Установка Docker..."
sudo apt install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

sudo usermod -aG docker $USER
newgrp docker
echo "Docker установлен и добавлен в группу пользователей."

# Установка Node.js и npm
echo "Установка Node.js и npm..."
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

echo "Node.js $(node -v) и npm $(npm -v) установлены"

# Установка protoc и генераторов
echo "Установка protobuf компилятора..."
sudo apt install -y protobuf-compiler libprotobuf-dev

# Установка protoc-gen-go без полной Go-среды
echo "Установка protoc-gen-go..."
PROTOC_GEN_GO_VERSION="v1.32.0"
# Тип процессора
ARCH=$(dpkg --print-architecture)
curl -sSL "https://github.com/protocolbuffers/protobuf-go/releases/download/$PROTOC_GEN_GO_VERSION/protoc-gen-go.$PROTOC_GEN_GO_VERSION.linux.$ARCH.tar.gz" | \ 
sudo tar -xz -C /usr/local/bin protoc-gen-go

# Установка protoc-gen-go-grpc
# Устанавливается v1.0.0
sudo apt install protoc-gen-go-grpc

# Установка protobuf-ts для фронтенда
echo "Установка protobuf-ts..."
sudo npm install -g @protobuf-ts/plugin

# Запуск сервиса
echo "Для применения изменения необходимо сделать: newgrp docker; reboot"
echo "Запуск сервиса:"
echo "Из родительской директории /gym выполните команду: make bundle-gen"
