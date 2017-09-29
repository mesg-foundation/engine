FROM node:8.5.0

RUN npm install -g yarn

COPY package.json .
COPY yarn.lock .

RUN yarn

COPY . .

CMD ["./start-listener.sh"]
