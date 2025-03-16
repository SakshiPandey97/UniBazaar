import React from "react";

const TeamMemberCard = ({ name, position, image }) => {
  return (
    <div className="max-w-xs bg-white rounded-lg shadow-lg overflow-hidden hover:shadow-2xl transition-transform transform hover:scale-105">
      <img className="w-full h-48 object-cover" src={image} alt={name} />
      <div className="p-4 text-center">
        <h3 className="text-xl font-semibold text-gray-800">{name}</h3>
        <p className="text-gray-600">{position}</p>
      </div>
    </div>
  );
};

export default TeamMemberCard;
