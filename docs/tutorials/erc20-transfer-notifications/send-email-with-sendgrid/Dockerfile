FROM node:10.12.0
WORKDIR /app
COPY ./package* ./
RUN npm install
COPY . .
CMD [ "node", "index.js" ]