import React, { useEffect, useState } from 'react';
import { Chart } from 'primereact/chart';
import { getValueHistory } from '../../api/service';
import './GraphComponent.scss';

interface ChartData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    fill: boolean;
    borderColor: string;
    tension: number;
  }[];
}

interface ValueLog {
  id: number;
  value1: number;
  value2: number;
  created_at: string;
}

const GraphComponent: React.FC = () => {
  const [lineChartData, setLineChartData] = useState<ChartData | null>(null);
  const [error, setError] = useState<string | null>(null);

  const lineChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: true,
        position: 'top' as const,
      },
    },
    scales: {
      x: {
        ticks: {
          autoSkip: true,
          maxTicksLimit: 10,
        }
      }
    }
  };

  const fetchHistory = async () => {
    try {
      const history: ValueLog[] = await getValueHistory();
      if (history && history.length > 0) {
        const labels = history.map(v => new Date(v.created_at).toLocaleTimeString()).reverse();
        const value1Data = history.map(v => v.value1).reverse();
        const value2Data = history.map(v => v.value2).reverse();

        const chartData: ChartData = {
          labels,
          datasets: [
            {
              label: 'Value 1',
              data: value1Data,
              fill: false,
              borderColor: '#42A5F5',
              tension: 0.4,
            },
            {
              label: 'Value 2',
              data: value2Data,
              fill: false,
              borderColor: '#FFA726',
              tension: 0.4,
            },
          ],
        };
        setLineChartData(chartData);
        setError(null);
      } else {
        setLineChartData(null);
      }
    } catch (err) {
      setError('Falha ao buscar o histórico de valores.');
      console.error(err);
    }
  };

  useEffect(() => {
    fetchHistory();
    const interval = setInterval(fetchHistory, 2000); // Atualiza a cada 2 segundos

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="graph-component">
      <div className="graph-bar-title">Histórico de Valores</div>
      <div className="graph-bar-content">
        {error && <p>{error}</p>}
        {lineChartData ? (
          <Chart type="line" data={lineChartData} options={lineChartOptions} />
        ) : (
          !error && <div>Carregando dados do gráfico...</div>
        )}
      </div>
    </div>
  );
};

export default React.memo(GraphComponent);
