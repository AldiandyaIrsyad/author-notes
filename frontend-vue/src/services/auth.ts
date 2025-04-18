import axios from 'axios';

const API_URL = 'http://localhost:8080/v1/auth'; //TODO: move to env file


interface LoginPayload {
  username: string; 
  password: string;
}

interface RegisterPayload {
  username: string; 
  email: string;
  password: string;
}


interface LoginResponse {
    token: string;
}

export const authService = {
  async login(payload: LoginPayload): Promise<LoginResponse> {
    try {
      const response = await axios.post<LoginResponse>(`${API_URL}/login`, payload);
      return response.data;
    } catch (error) {
      // Handle or re-throw error for the component to catch
      console.error('Login failed:', error);
      throw error; // Re-throwing allows components to handle UI updates
    }
  },

  async register(payload: RegisterPayload): Promise<any> { // Adjust 'any' to a specific response type if known
    try {
      const response = await axios.post(`${API_URL}/register`, payload);
      return response.data;
    } catch (error) {
      console.error('Registration failed:', error);
      throw error;
    }
  },
};
