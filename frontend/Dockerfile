# Етап 1: Збірка (Node.js)
FROM node:24-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Етап 2: Запуск (Nginx)
FROM nginx:alpine
# Копіюємо зібрані файли з 1-го етапу в папку Nginx
COPY --from=build /app/dist /usr/share/nginx/html

# Магічний конфіг, щоб React Router працював при оновленні сторінки
RUN echo 'server { \
    listen 80; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html index.htm; \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]