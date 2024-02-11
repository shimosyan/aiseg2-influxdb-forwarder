console.log('hello world');
console.log(process.env.AISEG2_IP_ADDRESS);

const main = () => {
  console.log(new Date().getTime());
};

setInterval(main, 1000);
