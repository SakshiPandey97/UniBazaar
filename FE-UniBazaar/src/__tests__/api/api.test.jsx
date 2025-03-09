// import { vi } from 'vitest';
// import axios from 'axios';
// import { 
//   userLoginAPI, 
//   userRegisterAPI, 
//   getAllUsersAPI, 
//   getAllProductsAPI, 
//   postProductAPI 
// } from '../../api/axios'; 

// vi.mock('axios');

// global.localStorage = {
//   setItem: vi.fn(),
//   getItem: vi.fn(),
//   removeItem: vi.fn(),
//   clear: vi.fn(),
// };

// describe('API functions', () => {
//   afterEach(() => {
//     vi.clearAllMocks(); 
//   });

//   it('should handle successful user login', async () => {
//     const mockResponse = { data: { userId: '12345' } };
//     axios.post.mockResolvedValue(mockResponse);

//     const userLoginObject = { username: 'testuser', password: 'password123' };
//     const userId = await userLoginAPI({ userLoginObject });

//     expect(userId).toBe('12345');
//     expect(localStorage.setItem).toHaveBeenCalledWith('userId', '12345');
//     expect(axios.post).toHaveBeenCalledWith('http://127.0.0.1:4000/login', userLoginObject);
//   });

//   it('should handle user login failure', async () => {
//     const mockError = new Error('Login failed');
//     axios.post.mockRejectedValue(mockError);

//     const userLoginObject = { username: 'wronguser', password: 'wrongpassword' };

//     await expect(userLoginAPI({ userLoginObject })).rejects.toThrowError('Login failed');
//   });

//   it('should handle successful user registration', async () => {
//     const mockResponse = { data: { success: true } };
//     axios.post.mockResolvedValue(mockResponse);

//     const userRegisterObject = { username: 'newuser', email: 'newuser@example.com' };
//     const data = await userRegisterAPI({ userRegisterObject });

//     expect(data).toEqual(mockResponse.data);
//     expect(axios.post).toHaveBeenCalledWith('http://127.0.0.1:4000/signup', userRegisterObject);
//   });

//   it('should handle user registration failure', async () => {
//     const mockError = new Error('Registration failed');
//     axios.post.mockRejectedValue(mockError);

//     const userRegisterObject = { username: 'existinguser', email: 'existing@example.com' };

//     await expect(userRegisterAPI({ userRegisterObject })).rejects.toThrowError('Registration failed');
//   });

//   // it('should fetch all users', async () => {
//   //   const mockResponse = { data: [{ id: '1', name: 'User1' }, { id: '2', name: 'User2' }] };
//   //   axios.get.mockResolvedValue(mockResponse);

//   //   const users = await getAllUsersAPI('1'); 

//   //   expect(users).toEqual(mockResponse.data);
//   //   expect(axios.get).toHaveBeenCalledWith('http://127.0.0.1:8080/users?exclude=1');
//   // });

//   // it('should handle error when fetching users', async () => {
//   //   const mockError = new Error('Error fetching users');
//   //   axios.get.mockRejectedValue(mockError);

//   //   await expect(getAllUsersAPI('1')).rejects.toThrowError('Error fetching users');
//   // });

//   it('should fetch all products', async () => {
//     const mockResponse = { data: [{ id: 'p1', name: 'Product 1' }, { id: 'p2', name: 'Product 2' }] };
//     axios.get.mockResolvedValue(mockResponse);

//     const products = await getAllProductsAPI();

//     expect(products).toEqual(mockResponse.data);
//     expect(axios.get).toHaveBeenCalledWith('https://unibazaar-products.azurewebsites.net/products');
//   });

//   it('should handle error when fetching products', async () => {
//     const mockError = new Error('Error fetching products');
//     axios.get.mockRejectedValue(mockError);

//     await expect(getAllProductsAPI()).rejects.toThrowError('Error fetching products');
//   });

//   it('should post a new product successfully', async () => {
//     const mockResponse = { data: { id: 'p1', name: 'New Product' } };
//     axios.post.mockResolvedValue(mockResponse);

//     const formData = { name: 'New Product', price: 10.99 };
//     const response = await postProductAPI(formData);

//     expect(response).toEqual(mockResponse.data);
//     expect(axios.post).toHaveBeenCalledWith('https://unibazaar-products.azurewebsites.net/products', formData);
//   });

//   it('should handle error when posting a new product', async () => {
//     const mockError = new Error('Error posting product');
//     axios.post.mockRejectedValue(mockError);

//     const formData = { name: 'New Product', price: 10.99 };

//     await expect(postProductAPI(formData)).rejects.toThrowError('Error posting product');
//   });
// });
