import React, { memo } from "react";
import { Handle, Position } from "reactflow";

interface CustomNodeData {
  label: string;
  isConnectable: boolean;
}

interface CustomNodeProps {
  data: CustomNodeData;
}

const CustomNodeTwoOne: React.FC<CustomNodeProps> = ({ data }) => {
  const { label, isConnectable } = data;

  return (
    <>
      <Handle
        type="target"
        id="a"
        position={Position.Left}
        style={{ background: "black" }}
        isConnectable={isConnectable}
      />

      <Handle
        type="target"
        position={Position.Right}
        id="b"
        style={{ background: "black" }}
        isConnectable={isConnectable}
      />

      <Handle
        type="source"
        position={Position.Bottom}
        id="c"
        style={{ background: "#784be8" }}
        isConnectable={isConnectable}
      />

      <div>{label}</div>
    </>
  );
};

export default memo(CustomNodeTwoOne);
