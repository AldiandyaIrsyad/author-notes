import axios from 'axios';

const API_URL = 'http://localhost:8080/v1/auth'; // Adjust if your backend runs elsewhere

// Define interfaces for the request payloads based on backend DTOs
// (Assuming backend expects email and password for login, and similar for register)
interface LoginPayload {
  username: string; // Changed from email to match backend LoginRequest DTO
  password: string;
}

interface RegisterPayload {
  username: string; // Adjust fields based on backend v1.RegisterRequest
  email: string;
  password: string;
}

// Define interface for the expected login response
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
