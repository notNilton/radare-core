// ReconciliationUtils.js

export const createAdjacencyMatrix = (nodes: any[], edges: any[]) => {
  const cnOneTwoNodes = nodes.filter((node) => node.type === "cnOneTwo");
  const incidencematrix = Array.from({ length: cnOneTwoNodes.length }, () =>
    Array(edges.length).fill(0)
  );

  edges.forEach((edge, edgeIndex) => {
    const sourceIndex = cnOneTwoNodes.findIndex(
      (node) => node.id === edge.source
    );
    const targetIndex = cnOneTwoNodes.findIndex(
      (node) => node.id === edge.target
    );

    if (sourceIndex !== -1) incidencematrix[sourceIndex][edgeIndex] = -1;
    if (targetIndex !== -1) incidencematrix[targetIndex][edgeIndex] = 1;
  });

  return incidencematrix;
};

export const calcularReconciliacao = (
  nodes: any[],
  edges: any[],
  reconciliarApi: (
    incidenceMatrix: any,
    values: any,
    tolerances: any,
    names: any[],
    atualizarProgresso: any
  ) => Promise<void>,
  atualizarProgresso: (message: string) => void,
  edgeNames: any[]
) => {
  const value = edges.map((edge) => edge.value);
  const tolerance = edges.map((edge) => edge.tolerance);
  const incidencematrix = createAdjacencyMatrix(nodes, edges);

  atualizarProgresso("Chamando API de reconciliação...");
  reconciliarApi(incidencematrix, value, tolerance, edgeNames, atualizarProgresso);
};

export const reconciliarApi = async (
  incidence_matrix: number[][],
  values: number[],
  tolerances: number[],
  names: string[],
  atualizarProgresso: (message: string) => void,
  jsonFile?: File
) => {
  try {
    atualizarProgresso("Enviando dados para o servidor...");

    let unreconciledata;

    if (jsonFile) {
      const fileContent = await jsonFile.text();
      const jsonData = JSON.parse(fileContent);
      unreconciledata = jsonData.unreconciledata || jsonData;
    } else {
      unreconciledata = [
        {
          values,
          tolerances,
        },
      ];
    }

    const timestamp = new Date().toISOString();

    const pacote = {
      data: {
        description: "Reconciliation for Q3 financial data across departments",
        user: "John Doe",
        timestamp,
        names,
        incidence_matrix,
        unreconciledata,
      },
    };

    console.log("Pacote a ser enviado:", pacote);

    const response = await fetch("http://localhost:5000/reconcile", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(pacote),
    });

    if (response.ok) {
      atualizarProgresso("Reconciliação bem-sucedida.");

      // Realiza o GET para substituir os dados no localStorage
      const getResponse = await fetch("http://localhost:5000/reconciled-data");

      if (getResponse.ok) {
        const reconciledDataArray = await getResponse.json();
        console.log("Dados reconciliados recebidos do servidor:", reconciledDataArray);

        // Formata os dados recebidos para o localStorage
        const formattedData = reconciledDataArray.map((entry: any) => ({
          id: entry[0],
          user: entry[1] || "postgres",
          time: entry[2] || timestamp,
          tagname: entry[3] || names,
          tagreconciled: entry[4]?.map((v: any) => (typeof v === 'number' ? v.toFixed(2) : v)) || values.map(v => v.toFixed(2)),
          tagcorrection: entry[5]?.map((v: any) => (typeof v === 'number' ? v.toFixed(2) : v)) || tolerances,
          tagmatrix: entry[6] || incidence_matrix,
        }));

        // Substitui completamente o localStorage com os dados obtidos
        localStorage.setItem("reconciliationData", JSON.stringify(formattedData));
        console.log("Dados atualizados no localStorage:", formattedData);

        // Dispara o evento para atualizar os componentes que dependem dos dados
        window.dispatchEvent(new CustomEvent("localStorageUpdated"));
      } else {
        console.error("Falha ao obter dados reconciliados.");
      }
    } else {
      console.error("Falha na reconciliação dos dados.");
      atualizarProgresso("Falha na reconciliação.");
    }
  } catch (error) {
    console.error("Erro ao reconciliar dados:", error);
    atualizarProgresso("Erro durante a reconciliação.");
  }
};
