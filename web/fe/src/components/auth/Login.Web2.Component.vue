<template>
  <q-layout view="lHh Lpr lFf">
    <q-page-container>
      <q-form @submit.prevent="onSubmit" class="q-gutter-md">
        <q-page class="flex flex-center bg-grey-2">
          <q-card class="q-pa-md shadow-2 my_card" bordered>
            <q-card-section class="text-center">
              <div class="text-grey-9 text-h5 text-weight-bold">Sign in</div>
              <div class="text-grey-8">Sign in below to access your account</div>
            </q-card-section>
            <q-card-section>
              <q-input v-model="email" dense outlined label="Email Address"></q-input>
              <q-input v-model="password" dense outlined class="q-mt-md" type="password" label="Password"></q-input>
            </q-card-section>
            <q-card-section>
              <q-btn style="border-radius: 8px" color="dark" rounded size="md" label="Sign in" type="submit" no-caps
                class="full-width"></q-btn>
            </q-card-section>
            <q-card-section class="text-center q-pt-none">
              <div class="text-grey-8">
                Don't have an account yet? Let contact with admin.
                <!-- <a href="#" class="text-dark text-weight-bold" style="text-decoration: none">Sign up.</a> -->
              </div>
            </q-card-section>
          </q-card>
        </q-page>
      </q-form>
    </q-page-container>
  </q-layout>
</template>

<script>

import loginService from 'src/services/auth/login.services'

export default {
  name: "FormAuth",

  data() {
    return {
      email: "",
      password: "",
    };
  },

  methods: {
    // Phương thức hiển thị thông báo
    showNotif(message, color) {
      this.$q.notify({
        message: "The website states: " + message,
        color: color,
        multiLine: true,
        avatar: "https://cdn.quasar.dev/img/boy-avatar.png",
        actions: [
          {
            label: "Close",
            color: "yellow",
            handler: () => {
              /* ... */
            },
          },
        ],
      });
    },

    // Phương thức hiển thị loading
    showLoading() {
      this.$q.loading.show({
        message: "Some important process is in progress. Hang on...",
      });
      // Ẩn loading sau 1 giây (1000ms)
      setTimeout(() => {
        this.$q.loading.hide();
      }, 1000);
    },

    // Phương thức xử lý form submit
    async onSubmit() {
      // const authStore = useAuthStore();
      let data = {
        email: "",
        password: "",
      };

      // Validate email
      if (this.email.length > 0) {
        const emailRegex = /^[\w-\\.]+@([\w-]+\.)+[\w-]{2,4}$/g;
        if (emailRegex.test(this.email)) {
          data.email = this.email;
        } else {
          this.showNotif(
            "The email format is incorrect. Please enter a valid email address.",
            "warning"
          );
          return;
        }
      } else {
        this.showNotif("Email field cannot be empty.", "warning");
        return;
      }

      // Validate password
      if (this.password.length <= 0) {
        this.showNotif("Password field cannot be empty.", "warning");
        return;
      }
      data.password = this.password;

      try {
        // const response = await AuthService.login(data);
        const response = await loginService.login(data);
        console.log(response);

        if (response.status_code != 200) {
          if (response.data.message == "Cannot login, wrong password account") {
            this.showNotif(
              "Incorrect email or password. Please check and try again.",
              "negative"
            );
            return;
          } else if (response.message == "Cannot login, wrong email account") {
            this.showNotif(
              "Email or password incorrect or not registered. Please contact the website administrator for more information.",
              "negative"
            );
            return;
          }

          throw response;
        } else {
          this.showLoading();
          this.showNotif(
            "Welcome to our website! We hope you have a great experience.",
            "positive"
          );
          // Redirect after successful login
          this.$router.push("/");
        }
      } catch (error) {
        this.showNotif("Login failed. Please try again later.", "negative");
        console.error("Login failed:", error);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.my_card {
  width: 25rem;
  border-radius: 8px;
  box-shadow: 0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
}
</style>
