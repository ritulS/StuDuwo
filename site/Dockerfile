FROM oven/bun AS build
WORKDIR /app
COPY bun.lockb .
COPY package.json .
RUN bun install --frozen-lockfile
COPY . .
RUN bun run build

FROM nginx:latest
WORKDIR /usr/share/nginx/html
EXPOSE 80
COPY --from=build /app/dist ./
CMD ["nginx", "-g", "daemon off;"]