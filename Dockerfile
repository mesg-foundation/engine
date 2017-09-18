FROM node:8.5.0

RUN npm install -g yarn

COPY package.json .

RUN yarn

COPY . .

CMD ["./start-listener.sh"]
