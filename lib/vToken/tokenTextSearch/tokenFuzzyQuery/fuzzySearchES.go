package tokenFuzzyQuery

import (
	"github.com/openGemini/openGemini/lib/utils"
	"github.com/openGemini/openGemini/lib/vToken/tokenIndex"
	"github.com/openGemini/openGemini/lib/vToken/tokenTextSearch/tokenMatchQuery"
	"sort"
)

func AbsInt(a int) int {
	if a >= 0 {
		return a
	} else {
		return -a
	}
}
func MinThree(a, b, c int) int {
	var min int
	if a >= b {
		min = b
	} else {
		min = a
	}
	if min >= c {
		return c
	} else {
		return min
	}
}
func minDistanceToken(word1 string, word2 string) int {
	l1 := len(word1)
	l2 := len(word2)

	dp := make([][]int, l1+1)

	dp[0] = make([]int, l2+1)
	for j := 0; j <= l2; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= l1; i++ {
		dp[i] = make([]int, l2+1)
		dp[i][0] = i
		for j := 1; j <= l2; j++ {
			if word1[i-1:i] == word2[j-1:j] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = MinThree(dp[i-1][j], dp[i][j-1], dp[i-1][j-1]) + 1
			}
		}
	}
	return dp[l1][l2]
}

func VerifyED(searStr string, dataStr string, distance int) bool {
	if minDistanceToken(searStr, dataStr) <= distance {
		return true
	} else {
		return false
	}
}
func UnionMapToken(map1 map[utils.SeriesId][]uint16, map2 map[utils.SeriesId][]uint16) map[utils.SeriesId][]uint16 {
	if len(map1) == 0 {
		return map2
	} else if len(map2) == 0 {
		return map1
	} else if len(map1) < len(map2) {
		for sid1, _ := range map1 {
			//if _,ok := map2[sid1]; !ok {
			map2[sid1] = make([]uint16, 0)
			//}
		}
		return map2
	} else {
		for sid1, _ := range map2 {
			//if _, ok := map1[sid1]; !ok {
			map1[sid1] = make([]uint16, 0)
			//}
		}
		return map1
	}
}
func UnionMapTokenThree(x map[utils.SeriesId][]uint16, y map[utils.SeriesId][]uint16, z map[utils.SeriesId][]uint16) map[utils.SeriesId][]uint16 {
	if len(x) == 0 {
		return UnionMapToken(y, z)
	} else if len(y) == 0 {
		return UnionMapToken(x, z)
	} else if len(z) == 0 {
		return UnionMapToken(x, y)
	} else {
		x = UnionMapToken(x, y)
		x = UnionMapToken(x, z)
		return x
	}
}
func ReadInvertedIndex(token string, indexRootNode *tokenIndex.IndexTreeNode) map[utils.SeriesId][]uint16 {
	var invertIndex tokenIndex.Inverted_index
	var indexNode *tokenIndex.IndexTreeNode
	var invertIndex2 tokenIndex.Inverted_index
	var invertIndex3 tokenIndex.Inverted_index

	tokenArr := []string{token}
	invertIndex, indexNode = tokenMatchQuery.SearchInvertedListFromCurrentNode(tokenArr, indexRootNode, 0, invertIndex, indexNode)
	invertIndexDeep := tokenMatchQuery.DeepCopy(invertIndex)
	invertIndex2 = tokenMatchQuery.SearchInvertedListFromChildrensOfCurrentNode(indexNode, nil)
	if indexNode != nil && len(indexNode.AddrOffset()) > 0 {
		invertIndex3 = tokenMatchQuery.TurnAddr2InvertLists(indexNode.AddrOffset(), invertIndex3)
	}
	//todo 三个一起合并应会快一些
	invertIndexDeep = UnionMapTokenThree(invertIndexDeep, invertIndex2, invertIndex3)
	return invertIndexDeep
}
func FuzzySearchComparedWithES(searchSingleToken string, indexRootNode *tokenIndex.IndexTreeNode, distance int) []utils.SeriesId {
	sum := 0
	sumPass := 0
	mapRes := make(map[utils.SeriesId][]uint16)
	q := 2
	lensearchToken := len(searchSingleToken)
	var qgramSearch = make([]FuzzyPrefixGram, 0)
	for i := 0; i < lensearchToken-q+1; i++ {
		qgramSearch = append(qgramSearch, NewFuzzyPrefixGram(searchSingleToken[i:i+q], i))
	}
	sort.SliceStable(qgramSearch, func(i, j int) bool {
		if qgramSearch[i].Gram() < qgramSearch[j].Gram() {
			return true
		}
		return false
	})
	prefixgramcount := q*distance + 1

	var mapsearchGram = make(map[string][]int)
	if lensearchToken-q+1 >= prefixgramcount {
		for i := 0; i < prefixgramcount; i++ {
			mapsearchGram[qgramSearch[i].Gram()] = append(mapsearchGram[qgramSearch[i].Gram()], qgramSearch[i].Pos())
		}
	}
	for i, _ := range indexRootNode.Children() {
		lenChildrendata := len(indexRootNode.Children()[i].Data())
		if lenChildrendata > lensearchToken+distance || lenChildrendata < lensearchToken-distance {
			continue
		} else if lenChildrendata-q+1 < prefixgramcount || lensearchToken-q+1 < prefixgramcount {
			sum++
			verifyresult := VerifyED(searchSingleToken, indexRootNode.Children()[i].Data(), distance)
			if verifyresult {
				sumPass++
				newMap := ReadInvertedIndex(indexRootNode.Children()[i].Data(), indexRootNode)
				mapRes = UnionMapToken(mapRes, newMap)
			}
			continue
		} else {
			flagCommon := 0
			var qgramData = make([]FuzzyPrefixGram, 0)
			for m := 0; m < lenChildrendata-q+1; m++ {
				qgramData = append(qgramData, NewFuzzyPrefixGram(indexRootNode.Children()[i].Data()[m:m+q], m))
			}
			sort.SliceStable(qgramData, func(m, n int) bool {
				if qgramData[m].Gram() < qgramData[n].Gram() {
					return true
				}
				return false
			})
			for k := 0; k < prefixgramcount; k++ {
				_, ok := mapsearchGram[qgramData[k].Gram()]
				if ok {
					for n := 0; n < len(mapsearchGram[qgramData[k].Gram()]); n++ {
						if AbsInt(mapsearchGram[qgramData[k].Gram()][n]-qgramData[k].Pos()) <= distance {
							flagCommon = 1
							sum++
							verifyresult2 := VerifyED(searchSingleToken, indexRootNode.Children()[i].Data(), distance)
							if verifyresult2 {
								sumPass++
								newMap := ReadInvertedIndex(indexRootNode.Children()[i].Data(), indexRootNode)
								mapRes = UnionMapToken(mapRes, newMap)

							}
							break
						}
					}
					if flagCommon == 1 {
						break
					}
				}
			}
			continue
		}
	}
	resArr := make([]utils.SeriesId, 0)
	for i, _ := range mapRes {
		resArr = append(resArr, i)
	}
	return resArr
}
