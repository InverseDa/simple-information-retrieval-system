# simple-information-retrieval-system

A simple information retrieval system for the 40th Anniversary of Shenzhen University.

## Framework

- Frontend: Vue3

- Backend: Gin

## How to deploy this project?

### 1. You need to install the following dependencies

```
go@1.19
nodejs@lts
python@3.10
pnpm
pip
```

For Ubuntu/Debian

```
sudo apt install golang-1.19 nodejs npm python3 pip
npm install -g pnpm
```

For Fedora/CentOS

```
sudo dnf install golang python3 python3-pip nodejs npm
npm install -g pnpm
```

For Arch/Manjaro

```
sudo pacman -S go nodejs pnpm python3 python3-pip
```

For MacOS

```
brew install golang nodejs pnpm python
```

### 2. Run Frontend

Direct to frontend project directory

```
cd ./frontend
```

Compiles and hot-reloads for development

```
pnpm run serve
```

Compiles and minifies for production

```
pnpm run build
```

Lints and fixes files

```
pnpm run lint
```

### 3. Run Backend

Direct to frontend project directory

```
cd ./backend
```

Install packages

```
go get
```

Now you can run Gin

```
go run main.go
```

## Where is config.ini?

`config.ini` includ your www1.szu.edu.cn account and password.

Format:

```config
[Profile]
account =
password =
```

It should be created in `backend/python` and loaded when `backend/python/web_scrapyer.py` is running.
