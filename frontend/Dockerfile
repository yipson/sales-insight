# Frontend Dockerfile
FROM node:20-alpine

WORKDIR /app

# Copiar package.json y package-lock.json
COPY package*.json ./
RUN npm install

# Copiar código fuente
COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host"]
