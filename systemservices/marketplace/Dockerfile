FROM node:10.11.0
WORKDIR /app
COPY ./package* ./
RUN npm install
COPY . .
CMD [ "npm", "start" ]