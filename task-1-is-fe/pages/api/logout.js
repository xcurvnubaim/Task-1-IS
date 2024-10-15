import cookie from 'cookie';

export default async function handler(req, res) {
  if (req.method === 'POST') {
    // Clear the auth-token cookie by setting it to expire in the past
    res.setHeader('Set-Cookie', cookie.serialize('auth-token', '', {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      expires: new Date(0), // Expire immediately
      path: '/',
    }));

    res.status(200).json({ message: 'Logout successful!' });
  } else {
    res.status(405).json({ message: 'Method not allowed' });
  }
}
