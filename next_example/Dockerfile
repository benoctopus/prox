FROM node:latest

EXPOSE 3000

WORKDIR /app

COPY . .

RUN npm install
RUN npm run build

CMD ["npm", "run", "start"]