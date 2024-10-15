import cookie from 'cookie';

export default async function handler(req, res) {
  if (req.method === 'POST') {
    const { username, password } = req.body;

    // Replace with your actual authentication logic
    if (username === 'test' && password === 'password') {
      const token = 'your-jwt-token'; // Replace with a real JWT token

      // Set the cookie
      res.setHeader('Set-Cookie', cookie.serialize('auth-token', token, {
        httpOnly: true, // Secure the cookie (not accessible via JavaScript)
        secure: process.env.NODE_ENV === 'production', // Use secure cookie in production
        maxAge: 60 * 60 * 24 * 7, // 1 week
        path: '/',
      }));

      res.status(200).json({ message: 'Login successful!' });
    } else {
      res.status(401).json({ message: 'Invalid credentials' });
    }
  } else {
    res.status(405).json({ message: 'Method not allowed' });
  }
}
