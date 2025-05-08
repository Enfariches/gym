# Подготовка
echo "Update archive.ubuntu"
sudo apt update && sudo apt upgrade -y

# Установка базовых утилит
echo "Installing make..."
sudo apt install -y ca-certificates curl gnupg make

# Установка Docker
echo "Installing Docker..."
sudo apt install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

sudo usermod -aG docker $USER

# Установка Node.js и npm
echo "Insalling Node.js..."
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# Установка protoc и генераторов
echo "Installing protobuf..."
sudo apt install -y protobuf-compiler libprotobuf-dev

# Установка protoc-gen-go без полной Go-среды
echo "Insatlling protoc-gen-go..."
sudo apt install protoc-gen-go-1-5

# Установка protoc-gen-go-grpc
# Устанавливается v1.0.0
echo "Insatlling protoc-gen-go-grpc..."
sudo apt install protoc-gen-go-grpc

# Установка protobuf-ts для фронтенда
echo "Установка protobuf-ts..."
sudo npm install -g @protobuf-ts/plugin

sudo newgrp docker

# Запуск сервиса
echo "Start Gymnastics..."
sudo make bundle-gen