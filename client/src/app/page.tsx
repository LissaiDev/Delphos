import Dashboard from '@/components/dashboard/Dashboard';
import { ENDPOINTS } from '@/config/constants';

export default function Home() {
  return (
    <main>
      <Dashboard endpoint={ENDPOINTS.DEFAULT_SSE} />
    </main>
  );
}