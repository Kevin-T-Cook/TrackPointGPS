<template>
  <div class="login-container">
    <img src="@/assets/TrackPoint.png" alt="TrackPoint Logo" class="login-logo" />
    <h1>Login</h1>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">Username</label>
        <input type="text" id="username" v-model="username" placeholder="Enter your username" required />
      </div>
      <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" v-model="password" placeholder="Enter your password" required />
      </div>
      <button type="submit" class="login-button">Login</button>
      <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
    </form>
  </div>
</template>

<script>
import api from "../api/api";

export default {
  name: "LoginForm",
  data() {
    return {
      username: "",
      password: "",
      errorMessage: null,
    };
  },
  methods: {
    async handleLogin() {
      try {
        const response = await api.login({
          username: this.username,
          password: this.password,
        });
        localStorage.setItem("token", response.token);
        this.$router.push("/");
      } catch (error) {
        this.errorMessage = "Login failed. Please check your credentials.";
        console.error("Login error:", error);
      }
    },
  },
};
</script>

<style scoped>
.login-container {
  max-width: 800px;
  margin: 50px auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 5px;
  background-color: #f9f9f9;
}

h1 {
  text-align: center;
  font-size: 2rem;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  font-weight: bold;
  margin-bottom: 5px;
}

input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 3px;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #f57c00;
  color: #fff;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 18px;
  font-weight: bold;
}

button:hover {
  background-color: #e65100;
}

.error-message {
  color: red;
  font-size: 14px;
  text-align: center;
  margin-top: 10px;
}

.login-logo {
  display: block;
  max-width: 100%;
  height: auto;
  margin: 0 auto 20px;
}

@media (max-width: 768px) {
  .login-container {
    max-width: 90%;
    padding: 15px;
  }

  h1 {
    font-size: 1.8rem;
  }

  button {
    font-size: 16px;
  }

  input {
    padding: 8px;
    font-size: 14px;
  }
}

@media (max-width: 480px) {
  .login-container {
    max-width: 95%;
    padding: 10px;
  }

  h1 {
    font-size: 1.5rem;
  }

  button {
    font-size: 14px;
    padding: 8px;
  }

  input {
    padding: 6px;
    font-size: 12px;
  }
}
</style>