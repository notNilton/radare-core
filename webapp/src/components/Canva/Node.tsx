import React, { useCallback, useState, useRef, useEffect } from "react";
import ReactFlow, {
  Controls,
  Panel,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
  BackgroundVariant,
  Connection,
  Edge,
  Node,
} from "reactflow";
import "reactflow/dist/style.css";
import "./Node.scss";

import {
  initialNodes,
  initialEdges,
  nodeTypes,
} from "./utils/initialCanvaDataIII";
import { calcularReconciliacao, reconciliarApi, createAdjacencyMatrix } from "./utils/Reconciliacao";
import SidebarComponent from "../Sidebar/SidebarComponent";
import GraphComponent from "../Graph/GraphComponent";
import PanelButtons from "./PanelButtons";

// Funções utilitárias para a funcionalidade de reconciliação (atualmente não integrada).
const generateRandomName = () => {
  const names = ["Laravel", "Alucard", "Sigma", "Delta", "Orion", "Phoenix"];
  return names[Math.floor(Math.random() * names.length)];
};
const getNodeId = () => `randomnode_${+new Date()}`;

interface EdgeData {
  nome: string;
  value?: number;
  tolerance?: number;
}

/**
 * NodeView é o componente principal que renderiza a interface do ReactFlow,
 * a barra lateral com os dados atuais e o gráfico com o histórico de dados.
 */
const NodeView: React.FC = () => {
  // Estados para gerenciar os nós e arestas do ReactFlow.
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState<EdgeData>(
    initialEdges.map((edge) => {
      const data = edge.data || {};
      return {
        ...edge,
        data: {
          nome: data.nome || generateRandomName(),
          value: data.value,
          tolerance: data.tolerance,
        },
        label: `Nome: ${data.nome || "Novo"}, Valor: ${data.value || "N/A"}, Tolerância: ${data.tolerance || "N/A"}`,
        type: "step",
      };
    })
  );

  // Refs para acessar os valores mais recentes de nós e arestas em callbacks.
  const nodesRef = useRef(nodes);
  const edgesRef = useRef(edges);
  useEffect(() => { nodesRef.current = nodes; }, [nodes]);
  useEffect(() => { edgesRef.current = edges; }, [edges]);

  // Estados para controlar a visibilidade da barra lateral e do gráfico.
  const [isSidebarVisible, setIsSidebarVisible] = useState(true);
  const [isGraphVisible, setIsGraphVisible] = useState(true);

  // Callback para adicionar uma nova aresta quando dois nós são conectados.
  const onConnect = useCallback(
    (params: Connection) => {
      const nome = generateRandomName();
      const newEdge: Edge<EdgeData> = {
        id: `e${params.source}-${params.target}-${nome}`,
        source: params.source!,
        target: params.target!,
        data: { nome, value: undefined, tolerance: undefined },
        label: `Nome: ${nome}, Valor: N/A, Tolerância: N/A`,
        type: "step",
      };
      setEdges((eds) => addEdge(newEdge, eds));
    },
    [setEdges]
  );

  // Callback para editar os dados de uma aresta com um duplo clique.
  const onEdgeDoubleClick = (_: React.MouseEvent, edge: Edge<EdgeData>) => {
    const valueStr = window.prompt("Digite um valor para a aresta:");
    const toleranceStr = window.prompt("Digite um valor para a tolerância:");

    if (valueStr && toleranceStr && !isNaN(Number(valueStr)) && !isNaN(Number(toleranceStr))) {
      const value = parseFloat(valueStr);
      const tolerance = parseFloat(toleranceStr);
      setEdges((prevEdges) =>
        prevEdges.map((e) => {
          if (e.id === edge.id) {
            const newData: EdgeData = { ...(e.data || { nome: generateRandomName() }), value, tolerance };
            return { ...e, data: newData, label: `Nome: ${newData.nome}, Valor: ${value}, Tolerância: ${tolerance}` };
          }
          return e;
        })
      );
    }
  };

  // Callback para adicionar um novo nó ao canva.
  const addNode = useCallback(
    (nodeType: string) => {
      const newNode: Node = {
        id: getNodeId(),
        type: nodeType,
        data: { label: "Simples", isConnectable: true },
        style: { background: "white", border: "2px solid black", padding: "3px", width: "100px" },
        position: { x: Math.random() * window.innerWidth, y: Math.random() * window.innerHeight },
      };
      setNodes((nds) => nds.concat(newNode));
    },
    [setNodes]
  );

  // As funções handleFileUpload e handleReconcile fazem parte da funcionalidade de "reconciliação"
  // que não está conectada ao backend Go atual e depende de uma API externa.
  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      const edgeNames = edges.map((edge) => (edge.data ? edge.data.nome : ""));
      const incidenceMatrix = createAdjacencyMatrix(nodes, edges);
      reconciliarApi(incidenceMatrix, [], [], edgeNames, (message) => console.log(message), file);
    }
  };

  const handleReconcile = () => {
    const edgeNames = edges.map((edge) => (edge.data ? edge.data.nome : ""));
    calcularReconciliacao(nodes, edges, reconciliarApi, (message) => console.log(message), edgeNames);
  };

  // Funções para alternar a visibilidade dos painéis.
  const toggleSidebar = () => setIsSidebarVisible(!isSidebarVisible);
  const toggleGraph = () => setIsGraphVisible(!isGraphVisible);

  // Função de depuração para exibir o estado atual dos nós e arestas.
  const showNodesAndEdges = () => {
    alert(`Nodes: ${JSON.stringify(nodesRef.current, null, 2)}\nEdges: ${JSON.stringify(edgesRef.current, null, 2)}`);
  };

  return (
    <div className={`node-container ${isSidebarVisible ? "" : "sidebar-hidden"}`}>
      <div className="reactflow-component">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          nodeTypes={nodeTypes}
          fitView
          onEdgeDoubleClick={onEdgeDoubleClick}
        >
          <Panel position="top-left" className="top-left-panel custom-panel">
            <PanelButtons
              addNode={addNode}
              showNodesAndEdges={showNodesAndEdges}
              toggleSidebar={toggleSidebar}
              toggleGraph={toggleGraph}
              handleReconcile={handleReconcile}
              handleFileUpload={handleFileUpload}
              isSidebarVisible={isSidebarVisible}
              isGraphVisible={isGraphVisible}
            />
          </Panel>
          <Controls />
          <Background variant={BackgroundVariant.Dots} gap={12} size={1} />
        </ReactFlow>

        {isSidebarVisible && (
          <div className="sidebar-component">
            <SidebarComponent />
          </div>
        )}
      </div>

      {isGraphVisible && (
        <div className="graph-component">
          <GraphComponent />
        </div>
      )}
    </div>
  );
};

export default NodeView;
