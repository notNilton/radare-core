import React, { memo } from "react";
import { Handle, Position } from "reactflow";

const customNodeOneThree = ({ data }) => {
  const { label, isConnectable } = data;

  return (
    <>
      <Handle
        type="source"
        id="a"
        position={Position.Left}
        style={{ background: "black" }}
        isConnectable={isConnectable}
      ></Handle>
      {/* <div>Custom Number Picker Node:</div>
      <strong>{data.number}</strong>
      <input
        className="nodrag"
        type="number"
        onChange={(e) => data.onChange(Number(e.target.value))}
        defaultValue={data.number}
      /> */}
      <Handle
        type="source"
        position={Position.Right}
        id="b"
        style={{ background: "black" }}
        isConnectable={isConnectable}
      ></Handle>

      <Handle
        type="target"
        position={Position.Top}
        id="c"
        style={{ background: "#784be8" }}
        isConnectable={isConnectable}
      ></Handle>

      <div>{`${label}`}</div>
    </>
  );
};

export default memo(customNodeOneThree);
