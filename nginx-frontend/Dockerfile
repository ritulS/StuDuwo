FROM node:alpine as builder
WORKDIR /app
COPY . .
RUN npm i
RUN npm run build

FROM nginx:latest
WORKDIR /usr/share/nginx/html
COPY --from=builder dist ./