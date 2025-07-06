import Dashboard from '@/components/dashboard/Dashboard';

export default function Home() {
  return (
    <main>
      <Dashboard endpoint="http://localhost:8080/api/stats/sse" />
    </main>
  );
}