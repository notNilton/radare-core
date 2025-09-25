// src/api/service.ts

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api';

/**
 * Busca os valores mais recentes do backend.
 */
export const getCurrentValues = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/current-values`);
    if (!response.ok) {
      // Se a resposta não for OK, mas for um 404, pode ser que o banco de dados ainda esteja vazio.
      if (response.status === 404) {
        return null; // Retorna nulo para que o frontend possa lidar com isso.
      }
      throw new Error(`A resposta da rede não foi ok: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Erro ao buscar os valores atuais:', error);
    throw error;
  }
};

/**
 * Busca o histórico de valores do backend.
 */
export const getValueHistory = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/values/history`);
    if (!response.ok) {
      throw new Error(`A resposta da rede não foi ok: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Erro ao buscar o histórico de valores:', error);
    throw error;
  }
};