/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// DefaultNormalizeScore generates a Normalize Score function that can normalize the
// scores from [0, max(scores)] to [0, maxPriority]. If reverse is set to true, it
// reverses the scores by subtracting it from maxPriority.
// Note: The input scores are always assumed to be non-negative integers.
func DefaultNormalizeScore(maxPriority int64, reverse bool, scores framework.NodeScoreList) *framework.Status {
	var maxCount int64
	// 寻找最大的分数, 节点分数中的最大值
	for i := range scores {
		if scores[i].Score > maxCount {
			maxCount = scores[i].Score
		}
	}

	// 什么情况下会出现这种情况呢？ 所有节点的分数都是 0
	if maxCount == 0 {
		if reverse {
			// 翻转的话，所有节点的分数都是 maxPriority
			// 默认都是 100 分
			for i := range scores {
				scores[i].Score = maxPriority
			}
		}
		return nil
	}

	for i := range scores {
		score := scores[i].Score

		// 重新计算分数, 归一化, 使得分数在 [0, maxPriority] 之间
		score = maxPriority * score / maxCount
		if reverse {
			// 翻转的话，原始分数越大，得分越低, 主要是给污点容忍插件使用的
			score = maxPriority - score
		}

		scores[i].Score = score
	}
	return nil
}
