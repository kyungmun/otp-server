const app = Vue.createApp({
  //delimiters: ['{{', '}}'],
  data() {
    return {
      title : 'OTP Server (TOTP), by kyungmun',
      imageLink : 'https://pkg.go.dev/static/shared/gopher/airplane-1200x945.svg',
      otpdatas : []
    };
  },
  created() {
    this.refrash()
  },
  methods : {
    refrash() {      
      axios.get('api/v1/otp')
        .then(response =>
          (this.otpdatas = response.data.data))
    }
  }
});

app.mount('#otplist');
