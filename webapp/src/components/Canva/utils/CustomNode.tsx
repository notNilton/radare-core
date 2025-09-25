import React, { memo } from "react";
import { Handle, Position } from "reactflow";

interface CustomNodeData {
  label: string;
  isConnectable: boolean;
  maxConnections?: number; // Tornando opcional, já que não é mais uma prop válida para Handle
}

interface CustomNodeProps {
  data: CustomNodeData;
}

const CustomNode: React.FC<CustomNodeProps> = ({ data }) => {
  const { label, isConnectable } = data;

  return (
    <>
      <Handle
        type="target"
        position={Position.Left}
        style={{ background: "black" }}
        isConnectable={isConnectable}
      />

      <Handle
        type="source"
        position={Position.Right}
        id="b"
        style={{ top: 20, background: "#784be8" }}
        isConnectable={isConnectable}
      />

      <Handle
        type="source"
        position={Position.Right}
        id="a"
        style={{ top: 75, background: "#784be8" }}
        isConnectable={isConnectable}
      />

      <div>{label}</div>
    </>
  );
};

export default memo(CustomNode);
