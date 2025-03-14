FROM node:latest
WORKDIR /app
COPY package.json /app/
RUN npm install
COPY ./writer.js /app/
CMD [ "node", "writer.js"]
