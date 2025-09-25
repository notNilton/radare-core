import React, { useState, useEffect } from 'react';
import { getCurrentValues } from '../../api/service';
import './SidebarComponent.scss';

interface CurrentValue {
  id: number;
  value1: number;
  value2: number;
  created_at: string;
}

const SidebarComponent: React.FC = () => {
  const [currentValue, setCurrentValue] = useState<CurrentValue | null>(null);
  const [error, setError] = useState<string | null>(null);

  const fetchCurrentValues = async () => {
    try {
      const data = await getCurrentValues();
      if (data) {
        setCurrentValue(data);
        setError(null);
      }
    } catch (err) {
      setError('Falha ao buscar os valores atuais.');
      console.error(err);
    }
  };

  useEffect(() => {
    fetchCurrentValues();
    const interval = setInterval(fetchCurrentValues, 2000); // Atualiza a cada 2 segundos

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="sidebar-container">
      <div className="sidebar-title">Valores Atuais</div>
      <div className="sidebar-content">
        {error && <p className="error-message">{error}</p>}
        {currentValue ? (
          <div className="values-display">
            <div className="value-item">
              <span className="value-label">Value 1:</span>
              <span className="value-data">{currentValue.value1}</span>
            </div>
            <div className="value-item">
              <span className="value-label">Value 2:</span>
              <span className="value-data">{currentValue.value2}</span>
            </div>
            <div className="timestamp">
              Última atualização: {new Date(currentValue.created_at).toLocaleString()}
            </div>
          </div>
        ) : (
          !error && <p>Carregando valores...</p>
        )}
      </div>
    </div>
  );
};

export default SidebarComponent;
