import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import Navbar from '@/components/layout/Navbar';
import Footer from '@/components/layout/Footer';
import ResumeUploadFloatingCard from '@/components/common/ResumeUploadFloatingCard';
import E2EStoreInit from '@/app/e2e-store-init';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: '面试吧AI面试平台',
  description: '大厂AI面试特训平台',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body className={inter.className}>
        <Navbar />
        <ResumeUploadFloatingCard />
        <E2EStoreInit />
        <main className="min-h-screen py-8">{children}</main>
        <Footer />
      </body>
    </html>
  );
}
