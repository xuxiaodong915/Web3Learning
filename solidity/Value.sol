pragma solidity >=0.4.22 <0.7.0;

contract Voting {
    // 存储候选人得票数的映射
    mapping(string => uint256) public votes;
    
    // 存储候选人列表（避免重复）
    string[] public candidates;
    
    // 记录候选人是否已存在
    mapping(string => bool) public candidateExists;
    
    // 事件定义
    event Voted(string candidate, uint256 newVoteCount);
    event CandidateAdded(string candidate);
    event VotesReset();

    // 投票函数 - 允许用户投票给某个候选人
    function vote(string memory candidate) public {
        require(bytes(candidate).length > 0, "候选人名称不能为空");
        require(candidateExists[candidate], "候选人不存在");
        
        votes[candidate]++;
        emit Voted(candidate, votes[candidate]);
    }
    
    // 添加候选人函数
    function addCandidate(string memory candidate) public {
        require(bytes(candidate).length > 0, "候选人名称不能为空");
        require(!candidateExists[candidate], "候选人已存在");
        
        candidates.push(candidate);
        candidateExists[candidate] = true;
        votes[candidate] = 0; // 初始化票数为0
        emit CandidateAdded(candidate);
    }

    // 获取某个候选人的得票数
    function getVotes(string memory candidate) public view returns (uint256) {
        return votes[candidate];
    }
    
    // 获取候选人总数
    function getCandidateCount() public view returns (uint256) {
        return candidates.length;
    }
    
    // 获取指定索引的候选人
    function getCandidate(uint256 index) public view returns (string memory) {
        require(index < candidates.length, "索引超出范围");
        return candidates[index];
    }

    // 重置所有候选人的得票数
    function resetVotes() public {
        uint256 len = candidates.length;
        for (uint256 i = 0; i < len; i++) {
            votes[candidates[i]] = 0;
        }
        emit VotesReset();
    }
}