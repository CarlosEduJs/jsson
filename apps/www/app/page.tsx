import { Hero } from "@/components/landing/hero";
import { Features } from "@/components/landing/features";
import { CodeDemo } from "@/components/landing/code-demo";
import { FAQ } from "@/components/landing/faq";
import { Footer } from "@/components/landing/footer";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col bg-background">
      <Hero />
      <Features />
      <CodeDemo />
      <FAQ />
      <Footer />
    </main>
  );
}
