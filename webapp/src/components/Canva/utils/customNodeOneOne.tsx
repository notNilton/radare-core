import React, { memo } from "react";
import { Handle, Position } from "reactflow";

interface CustomNodeData {
  label: string;
  isConnectable: boolean;
}

interface CustomNodeProps {
  data: CustomNodeData;
}

const CustomNodeOneOne: React.FC<CustomNodeProps> = ({ data }) => {
  return (
    <>
      <Handle
        type="target"
        id="target"
        position={Position.Left}
        style={{ background: "black" }}
        isConnectable={data.isConnectable}
      />
      <Handle
        type="source"
        position={Position.Right}
        id="a"
        style={{ top: 50, background: "#784be8" }}
        isConnectable={data.isConnectable}
      />
      <div>{data.label}</div>
    </>
  );
};

export default memo(CustomNodeOneOne);
