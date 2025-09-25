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
import PanelButtons from "./PanelButtons"; // Importando o novo componente

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

const NodeView: React.FC = () => {
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

  const nodesRef = useRef(nodes);
  const edgesRef = useRef(edges);
  const [isSidebarVisible, setIsSidebarVisible] = useState(true);
  const [isGraphVisible, setIsGraphVisible] = useState(true);

  useEffect(() => {
    nodesRef.current = nodes;
  }, [nodes]);

  useEffect(() => {
    edgesRef.current = edges;
  }, [edges]);

  const onConnect = useCallback(
    (params: Connection) => {
      const nome = generateRandomName();
      const newEdge: Edge<EdgeData> = {
        id: `e${params.source}-${params.target}-${nome}`,
        source: params.source!,
        target: params.target!,
        data: {
          nome,
          value: undefined,
          tolerance: undefined,
        },
        label: `Nome: ${nome}, Valor: N/A, Tolerância: N/A`,
        type: "step",
      };
      setEdges((eds) => addEdge(newEdge, eds));
    },
    [setEdges]
  );

  const onEdgeDoubleClick = (_: React.MouseEvent, edge: Edge<EdgeData>) => {
    const valueStr = window.prompt("Digite um valor para a aresta:");
    const toleranceStr = window.prompt("Digite um valor para a tolerância:");

    if (valueStr && toleranceStr && !isNaN(Number(valueStr)) && !isNaN(Number(toleranceStr))) {
      const value = parseFloat(valueStr);
      const tolerance = parseFloat(toleranceStr);

      setEdges((prevEdges) =>
        prevEdges.map((e) => {
          if (e.id === edge.id) {
            const newData: EdgeData = {
              ...(e.data || { nome: generateRandomName() }),
              value,
              tolerance,
            };
            return {
              ...e,
              data: newData,
              label: `Nome: ${newData.nome}, Valor: ${value}, Tolerância: ${tolerance}`,
            };
          }
          return e;
        })
      );
    }
  };

  const addNode = useCallback(
    (nodeType: string) => {
      const newNode: Node = {
        id: getNodeId(),
        type: nodeType,
        data: { label: "Simples", isConnectable: true },
        style: {
          background: "white",
          border: "2px solid black",
          padding: "3px",
          width: "100px",
        },
        position: {
          x: Math.random() * window.innerWidth,
          y: Math.random() * window.innerHeight,
        },
      };
      setNodes((nds) => nds.concat(newNode));
    },
    [setNodes]
  );

  const handleFileUpload = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (file) {
      const edgeNames = edges.map((edge) => (edge.data ? edge.data.nome : ""));
      const incidenceMatrix = createAdjacencyMatrix(nodes, edges);
  
      reconciliarApi(
        incidenceMatrix,
        [],
        [],
        edgeNames,
        (message) => console.log(message),
        file
      );
    }
  };
  
  const handleReconcile = () => {
    const edgeNames = edges.map((edge) => (edge.data ? edge.data.nome : ""));
    calcularReconciliacao(
      nodes,
      edges,
      reconciliarApi,
      (message) => {
        console.log(message);
      },
      edgeNames
    );
  };
  
  const toggleSidebar = () => {
    setIsSidebarVisible(!isSidebarVisible);
  };

  const toggleGraph = () => {
    setIsGraphVisible(!isGraphVisible);
  };

  const showNodesAndEdges = () => {
    console.log("Nodes:", nodesRef.current);
    console.log("Edges:", edgesRef.current);
    alert(
      `Nodes: ${JSON.stringify(
        nodesRef.current,
        null,
        2
      )}\nEdges: ${JSON.stringify(edgesRef.current, null, 2)}`
    );
  };

  return (
    <div
      className={`node-container ${isSidebarVisible ? "" : "sidebar-hidden"}`}
    >
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
