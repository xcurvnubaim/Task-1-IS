import { redirect } from 'next/navigation';

export default function Home() {
  // Perform the redirect
  redirect('/dashboard');
}